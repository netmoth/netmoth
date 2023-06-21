package connection

import (
	"bytes"

	"github.com/google/gopacket"

	"github.com/netmoth/netmoth/internal/utils"
)

// NewUDP is ...
func NewUDP(packet gopacket.Packet, ci gopacket.CaptureInfo) *Connection {
	transportFlow := packet.TransportLayer().TransportFlow()
	networkFlow := packet.NetworkLayer().NetworkFlow()
	srcPort, dstPort, _ := utils.ProcessPorts(transportFlow)
	return &Connection{
		Timestamp:       ci.Timestamp,
		UID:             networkFlow.FastHash() + transportFlow.FastHash(),
		SourceIP:        networkFlow.Src().String(),
		SourcePort:      srcPort,
		DestinationIP:   networkFlow.Dst().String(),
		DestinationPort: dstPort,
		TransportType:   "udp",
		Payload:         bytes.NewBuffer(packet.TransportLayer().LayerPayload()),
		Analyzers:       make(map[string]interface{}),
	}
}
