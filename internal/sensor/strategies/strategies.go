package strategies

import (
	"github.com/google/gopacket"
	"github.com/netmoth/netmoth/internal/config"
)

const (
	clusterID = 1234
)

// PacketsCaptureStrategy is ...
type PacketsCaptureStrategy interface {
	New(c *config.Config) ([]PacketDataSource, error)
	Destroy()
	PacketStats() (received uint64, dropped uint64)
}

// PacketDataSource is ...
type PacketDataSource interface {
	gopacket.PacketDataSource
	gopacket.ZeroCopyPacketDataSource
}

// Strategies is ...
func Strategies() map[string]PacketsCaptureStrategy {
	return map[string]PacketsCaptureStrategy{
		"pcap":     &PCAPStrategy{},
		"pfring":   &PFringsStrategy{},
		"afpacket": &AFPacketStrategy{},
	}
}
