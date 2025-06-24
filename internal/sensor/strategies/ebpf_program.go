package strategies

import (
	"encoding/binary"
	"fmt"

	"github.com/cilium/ebpf"
)

// PacketEvent represents a packet event sent from eBPF to userspace
type PacketEvent struct {
	Timestamp uint64
	Length    uint32
	Data      [1500]byte // Max packet size
}

// createRingBufferMap creates a ring buffer map for packet events
func createRingBufferMap() (*ebpf.Map, error) {
	ringbufMap, err := ebpf.NewMap(&ebpf.MapSpec{
		Type:       ebpf.RingBuf,
		MaxEntries: 1024 * 1024, // 1MB ring buffer
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ring buffer map: %w", err)
	}

	return ringbufMap, nil
}

// createStatsMap creates a map for packet statistics
func createStatsMap() (*ebpf.Map, error) {
	statsMap, err := ebpf.NewMap(&ebpf.MapSpec{
		Type:       ebpf.Array,
		KeySize:    4, // uint32
		ValueSize:  8, // uint64
		MaxEntries: 2, // received, dropped
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create stats map: %w", err)
	}

	// Initialize stats
	received := uint64(0)
	dropped := uint64(0)

	if err := statsMap.Update(uint32(0), received, ebpf.UpdateAny); err != nil {
		return nil, fmt.Errorf("failed to initialize received stats: %w", err)
	}

	if err := statsMap.Update(uint32(1), dropped, ebpf.UpdateAny); err != nil {
		return nil, fmt.Errorf("failed to initialize dropped stats: %w", err)
	}

	return statsMap, nil
}

// parsePacketEvent parses a packet event from the ring buffer
func parsePacketEvent(data []byte) (*PacketEvent, error) {
	if len(data) < 12 { // timestamp + length + minimum data
		return nil, fmt.Errorf("packet event too short")
	}

	event := &PacketEvent{}

	// Parse timestamp (8 bytes)
	event.Timestamp = binary.LittleEndian.Uint64(data[0:8])

	// Parse length (4 bytes)
	event.Length = binary.LittleEndian.Uint32(data[8:12])

	// Copy packet data
	copy(event.Data[:], data[12:])

	return event, nil
}

// createPacketData creates packet data from the event
func createPacketData(event *PacketEvent) []byte {
	// Create a proper Ethernet frame
	packet := make([]byte, event.Length)

	// Copy the packet data
	copy(packet, event.Data[:event.Length])

	return packet
}

// updateStats updates packet statistics
func updateStats(statsMap *ebpf.Map, received, dropped uint64) error {
	if err := statsMap.Update(uint32(0), received, ebpf.UpdateAny); err != nil {
		return fmt.Errorf("failed to update received stats: %w", err)
	}

	if err := statsMap.Update(uint32(1), dropped, ebpf.UpdateAny); err != nil {
		return fmt.Errorf("failed to update dropped stats: %w", err)
	}

	return nil
}

// getStats retrieves packet statistics from the eBPF map
func getStats(statsMap *ebpf.Map) (received uint64, dropped uint64, err error) {
	var receivedVal, droppedVal uint64

	if err := statsMap.Lookup(uint32(0), &receivedVal); err != nil {
		return 0, 0, fmt.Errorf("failed to get received stats: %w", err)
	}

	if err := statsMap.Lookup(uint32(1), &droppedVal); err != nil {
		return 0, 0, fmt.Errorf("failed to get dropped stats: %w", err)
	}

	return receivedVal, droppedVal, nil
}
