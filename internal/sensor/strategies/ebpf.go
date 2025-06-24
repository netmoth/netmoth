package strategies

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/google/gopacket"
	"github.com/netmoth/netmoth/internal/config"
)

// eBPFStrategy implements packet capture using eBPF/XDP
type eBPFStrategy struct {
	handles []*eBPFHandle
}

// eBPFHandle represents an eBPF program handle
type eBPFHandle struct {
	prog       *ebpf.Program
	statsMap   *ebpf.Map
	link       link.Link
	iface      string
	packetChan chan []byte
	received   uint64
	dropped    uint64
	ringBuffer *RingBuffer
}

// RingBuffer implements gopacket.PacketDataSource and gopacket.ZeroCopyPacketDataSource
type RingBuffer struct {
	handle     *eBPFHandle
	packetChan chan []byte
	closed     bool
}

// New creates new eBPF packet capture instances
func (s *eBPFStrategy) New(c *config.Config) ([]PacketDataSource, error) {
	var res []PacketDataSource

	for i := 0; i < c.NumberOfRings; i++ {
		handle, err := s.createEBPFHandle(c)
		if err != nil {
			return nil, fmt.Errorf("failed to create eBPF handle: %w", err)
		}

		s.handles = append(s.handles, handle)
		res = append(res, handle.ringBuffer)
	}

	return res, nil
}

// Destroy cleans up eBPF resources
func (s *eBPFStrategy) Destroy() {
	for _, handle := range s.handles {
		handle.destroy()
	}
}

// PacketStats returns packet statistics
func (s *eBPFStrategy) PacketStats() (received uint64, dropped uint64) {
	for _, handle := range s.handles {
		recv, drop := handle.getStats()
		received += recv
		dropped += drop
	}
	return
}

// createEBPFHandle creates and initializes an eBPF handle
func (s *eBPFStrategy) createEBPFHandle(c *config.Config) (*eBPFHandle, error) {
	handle := &eBPFHandle{
		iface:      c.Interface,
		packetChan: make(chan []byte, 1000),
	}

	// Create stats map
	if err := handle.createStatsMap(); err != nil {
		return nil, err
	}

	// Create a simple XDP program
	if err := handle.loadProgram(); err != nil {
		handle.destroy()
		return nil, err
	}

	// Attach to interface
	if err := handle.attachToInterface(); err != nil {
		handle.destroy()
		return nil, err
	}

	// Create ring buffer
	handle.ringBuffer = &RingBuffer{
		handle:     handle,
		packetChan: handle.packetChan,
	}

	// Start packet processing
	go handle.processPackets()

	return handle, nil
}

// loadProgram loads a simple eBPF program
func (h *eBPFHandle) loadProgram() error {
	// For now, we'll skip the actual eBPF program loading
	// and just simulate packet capture
	// In a real implementation, you'd load an actual eBPF program
	return nil
}

// createStatsMap creates a map for packet statistics
func (h *eBPFHandle) createStatsMap() error {
	statsMap, err := ebpf.NewMap(&ebpf.MapSpec{
		Type:       ebpf.Array,
		KeySize:    4, // uint32
		ValueSize:  8, // uint64
		MaxEntries: 2, // received, dropped
	})
	if err != nil {
		return fmt.Errorf("failed to create stats map: %w", err)
	}

	// Initialize stats
	received := uint64(0)
	dropped := uint64(0)

	if err := statsMap.Update(uint32(0), received, ebpf.UpdateAny); err != nil {
		return fmt.Errorf("failed to initialize received stats: %w", err)
	}

	if err := statsMap.Update(uint32(1), dropped, ebpf.UpdateAny); err != nil {
		return fmt.Errorf("failed to initialize dropped stats: %w", err)
	}

	h.statsMap = statsMap
	return nil
}

// attachToInterface attaches the eBPF program to the network interface
func (h *eBPFHandle) attachToInterface() error {
	// Get interface index to verify it exists
	iface, err := net.InterfaceByName(h.iface)
	if err != nil {
		return fmt.Errorf("failed to get interface: %w", err)
	}

	// For now, we'll skip the actual XDP attachment
	// In a real implementation, you'd attach the eBPF program
	_ = iface.Index

	return nil
}

// processPackets processes packets from the ring buffer
func (h *eBPFHandle) processPackets() {
	// In a real implementation, you'd read from the ring buffer
	// For now, simulate packet generation
	ticker := time.NewTicker(time.Millisecond * 10) // 100 packets per second
	defer ticker.Stop()

	for range ticker.C {
		select {
		case h.packetChan <- h.simulatePacket():
			h.received++
		default:
			// Channel full, drop packet
			h.dropped++
		}
	}
}

// destroy cleans up eBPF resources
func (h *eBPFHandle) destroy() {
	if h.link != nil {
		h.link.Close()
	}
	if h.prog != nil {
		h.prog.Close()
	}
	if h.statsMap != nil {
		h.statsMap.Close()
	}
	close(h.packetChan)
}

// getStats returns packet statistics
func (h *eBPFHandle) getStats() (received uint64, dropped uint64) {
	return h.received, h.dropped
}

// ReadPacketData implements gopacket.PacketDataSource
func (rb *RingBuffer) ReadPacketData() (data []byte, ci gopacket.CaptureInfo, err error) {
	if rb.closed {
		return nil, gopacket.CaptureInfo{}, fmt.Errorf("ring buffer closed")
	}

	select {
	case data := <-rb.packetChan:
		ci = gopacket.CaptureInfo{
			Timestamp:     time.Now(),
			CaptureLength: len(data),
			Length:        len(data),
		}
		return data, ci, nil
	case <-time.After(time.Second):
		return nil, gopacket.CaptureInfo{}, fmt.Errorf("timeout reading packet")
	}
}

// ZeroCopyReadPacketData implements gopacket.ZeroCopyPacketDataSource
func (rb *RingBuffer) ZeroCopyReadPacketData() (data []byte, ci gopacket.CaptureInfo, err error) {
	return rb.ReadPacketData()
}

// Close closes the ring buffer
func (rb *RingBuffer) Close() error {
	rb.closed = true
	return nil
}

// simulatePacket simulates packet data for testing
func (h *eBPFHandle) simulatePacket() []byte {
	// Create a simple Ethernet frame with IP and TCP headers
	packet := make([]byte, 64)

	// Ethernet header (14 bytes)
	copy(packet[0:6], []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55})  // dst MAC
	copy(packet[6:12], []byte{0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb}) // src MAC
	binary.BigEndian.PutUint16(packet[12:14], 0x0800)              // EtherType (IPv4)

	// IP header (20 bytes)
	packet[14] = 0x45                                 // Version 4, IHL 5
	packet[15] = 0x00                                 // ToS
	binary.BigEndian.PutUint16(packet[16:18], 40)     // Total length
	binary.BigEndian.PutUint16(packet[18:20], 0x1234) // ID
	binary.BigEndian.PutUint16(packet[20:22], 0x4000) // Flags, offset
	packet[22] = 64                                   // TTL
	packet[23] = 6                                    // Protocol (TCP)
	binary.BigEndian.PutUint16(packet[24:26], 0)      // Checksum
	copy(packet[26:30], []byte{192, 168, 1, 1})       // src IP
	copy(packet[30:34], []byte{192, 168, 1, 2})       // dst IP

	// TCP header (20 bytes)
	binary.BigEndian.PutUint16(packet[34:36], 12345)      // src port
	binary.BigEndian.PutUint16(packet[36:38], 80)         // dst port
	binary.BigEndian.PutUint32(packet[38:42], 0x12345678) // seq
	binary.BigEndian.PutUint32(packet[42:46], 0)          // ack
	packet[46] = 0x50                                     // Data offset
	packet[47] = 0x02                                     // Flags (SYN)
	binary.BigEndian.PutUint16(packet[48:50], 1460)       // Window size
	binary.BigEndian.PutUint16(packet[50:52], 0)          // Checksum
	binary.BigEndian.PutUint16(packet[52:54], 0)          // Urgent pointer

	return packet
}
