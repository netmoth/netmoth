package connection

import (
	"github.com/google/gopacket"

	"github.com/netmoth/netmoth/internal/utils"
)

// NewUDP is ...
func NewUDP(packet gopacket.Packet, ci gopacket.CaptureInfo) *Connection {
	transportFlow := packet.TransportLayer().TransportFlow()
	networkFlow := packet.NetworkLayer().NetworkFlow()
	srcPort, dstPort, _ := utils.ProcessPorts(transportFlow)

	// Используем пул объектов
	conn := GlobalConnectionPool.Get()

	conn.Timestamp = ci.Timestamp
	conn.UID = networkFlow.FastHash() + transportFlow.FastHash()
	conn.SourceIP = networkFlow.Src().String()
	conn.SourcePort = srcPort
	conn.DestinationIP = networkFlow.Dst().String()
	conn.DestinationPort = dstPort
	conn.TransportType = "udp"

	// Копируем payload без дополнительных аллокаций
	payload := packet.TransportLayer().LayerPayload()
	if len(payload) > 0 {
		conn.Payload.Write(payload)
	}

	return conn
}
