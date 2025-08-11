package sensor

import (
	"testing"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/netmoth/netmoth/internal/connection"
)

func buildUDPPacket() gopacket.Packet {
	eth := &layers.Ethernet{SrcMAC: []byte{1, 2, 3, 4, 5, 6}, DstMAC: []byte{6, 5, 4, 3, 2, 1}, EthernetType: layers.EthernetTypeIPv4}
	ip := &layers.IPv4{Version: 4, IHL: 5, TTL: 64, Protocol: layers.IPProtocolUDP}
	ip.SrcIP = []byte{192, 168, 0, 1}
	ip.DstIP = []byte{192, 168, 0, 2}
	udp := &layers.UDP{SrcPort: 12345, DstPort: 80}
	_ = udp.SetNetworkLayerForChecksum(ip)
	payload := gopacket.Payload([]byte("ping"))
	buf := gopacket.NewSerializeBuffer()
	_ = gopacket.SerializeLayers(buf, gopacket.SerializeOptions{ComputeChecksums: true, FixLengths: true}, eth, ip, udp, payload)
	pkt := gopacket.NewPacket(buf.Bytes(), layers.LinkTypeEthernet, gopacket.DecodeStreamsAsDatagrams)
	m := pkt.Metadata()
	m.CaptureInfo = gopacket.CaptureInfo{Timestamp: time.Now(), CaptureLength: len(buf.Bytes()), Length: len(buf.Bytes())}
	return pkt
}

func buildTCPPacket() gopacket.Packet {
	eth := &layers.Ethernet{SrcMAC: []byte{1, 2, 3, 4, 5, 6}, DstMAC: []byte{6, 5, 4, 3, 2, 1}, EthernetType: layers.EthernetTypeIPv4}
	ip := &layers.IPv4{Version: 4, IHL: 5, TTL: 64, Protocol: layers.IPProtocolTCP}
	ip.SrcIP = []byte{192, 168, 0, 1}
	ip.DstIP = []byte{192, 168, 0, 2}
	tcp := &layers.TCP{SrcPort: 12345, DstPort: 80, SYN: true}
	_ = tcp.SetNetworkLayerForChecksum(ip)
	buf := gopacket.NewSerializeBuffer()
	_ = gopacket.SerializeLayers(buf, gopacket.SerializeOptions{ComputeChecksums: true, FixLengths: true}, eth, ip, tcp)
	pkt := gopacket.NewPacket(buf.Bytes(), layers.LinkTypeEthernet, gopacket.DecodeStreamsAsDatagrams)
	m := pkt.Metadata()
	m.CaptureInfo = gopacket.CaptureInfo{Timestamp: time.Now(), CaptureLength: len(buf.Bytes()), Length: len(buf.Bytes())}
	return pkt
}

func TestProcessPacket_UDP_Enqueue(t *testing.T) {
	s := &sensor{connections: make(chan *connection.Connection, 1)}
	pkt := buildUDPPacket()
	s.processPacket(pkt)
	select {
	case c := <-s.connections:
		if c.TransportType != "udp" {
			t.Fatalf("expected udp, got %s", c.TransportType)
		}
		if c.SourceIP == "" || c.DestinationIP == "" || c.SourcePort == 0 || c.DestinationPort == 0 {
			t.Fatalf("unexpected connection fields: %+v", c)
		}
		if c.Payload.Len() == 0 {
			t.Fatalf("expected payload")
		}
	case <-time.After(400 * time.Millisecond):
		t.Fatalf("udp connection not enqueued")
	}
}

func TestProcessPacket_TCP_Paths(t *testing.T) {
	s := &sensor{connections: make(chan *connection.Connection, 1)}
	pkt := buildTCPPacket()
	s.streamFactory = &connection.TCPStreamFactory{Connections: s.connections}
	s.streamFactory.CreateAssembler()
	s.streamFactory.Ticker = time.NewTicker(10 * time.Millisecond)
	defer s.streamFactory.Ticker.Stop()
	s.processPacket(pkt)
}

func TestAgentBufferLimit(t *testing.T) {
	s := &sensor{agentMode: true}
	// add more than 100000 connections
	for i := 0; i < 100000+1000; i++ {
		c := &connection.Connection{}
		s.addToBuffer(c)
	}
	if s.connSize > 100000 {
		t.Fatalf("buffer exceeded limit: %d", s.connSize)
	}
}
