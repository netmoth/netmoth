package strategies

import (
	"log"

	"github.com/google/gopacket/pcap"

	"github.com/netmoth/netmoth/internal/config"
)

// PCAPStrategy is ...
type PCAPStrategy struct {
	handle *pcap.Handle
}

// New  is ...
func (s *PCAPStrategy) New(c *config.Config) ([]PacketDataSource, error) {
	if c.NumberOfRings > 1 {
		log.Println("WARNING: pcap not support cluster mode, ignoring number of rings parameter")
	}

	var err error
	s.handle, err = pcap.OpenLive(c.Interface, int32(c.SnapLen), c.Promiscuous, pcap.BlockForever)
	if err != nil {
		return nil, err
	}
	if len(c.Bpf) != 0 {
		if err = s.handle.SetBPFFilter(c.Bpf); err != nil {
			return nil, err
		}
	}

	return []PacketDataSource{s.handle}, nil
}

// Destroy is ...
func (s *PCAPStrategy) Destroy() {
	s.handle.Close()
}

// PacketStats is ...
func (s *PCAPStrategy) PacketStats() (received uint64, dropped uint64) {
	stats, err := s.handle.Stats()
	if err != nil {
		return
	}

	if stats.PacketsReceived > 0 {
		received += uint64(stats.PacketsReceived)
	}

	if stats.PacketsDropped > 0 {
		dropped += uint64(stats.PacketsDropped)
	}
	return
}
