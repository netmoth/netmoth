package strategies

import (
	"github.com/google/gopacket/pfring"

	"github.com/netmoth/netmoth/internal/config"
)

// PFringsStrategy is ...
type PFringsStrategy struct {
	rings []*pfring.Ring
}

// New is ...
func (s *PFringsStrategy) New(c *config.Config) ([]PacketDataSource, error) {
	var res []PacketDataSource
	for i := 0; i < c.NumberOfRings; i++ {
		ring, err := pfring.NewRing(c.Interface, uint32(c.SnapLen), pfring.FlagPromisc)
		if err != nil {
			return nil, err
		}

		if err = ring.SetDirection(pfring.ReceiveAndTransmit); err != nil {
			return nil, err
		}
		if err = ring.SetSocketMode(pfring.ReadOnly); err != nil {
			return nil, err
		}
		if c.NumberOfRings > 1 {
			if err = ring.SetCluster(clusterID, pfring.ClusterPerFlow5Tuple); err != nil {
				return nil, err
			}
		}
		if len(c.Bpf) != 0 {
			if err = ring.SetBPFFilter(c.Bpf); err != nil {
				return nil, err
			}
		}
		if err = ring.Enable(); err != nil {
			return nil, err
		}

		s.rings = append(s.rings, ring)
		res = append(res, ring)
	}
	return res, nil
}

// Destroy is ...
func (s *PFringsStrategy) Destroy() {
	for _, ring := range s.rings {
		_ = ring.Disable()
	}
	for _, ring := range s.rings {
		ring.Close()
	}
}

// PacketStats is ...
func (s *PFringsStrategy) PacketStats() (received uint64, dropped uint64) {
	for _, ring := range s.rings {
		stats, err := ring.Stats()
		if err != nil {
			continue
		}
		received += stats.Received
		dropped += stats.Dropped
	}
	return
}
