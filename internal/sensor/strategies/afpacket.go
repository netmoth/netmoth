package strategies

import (
	"github.com/google/gopacket/afpacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"golang.org/x/net/bpf"

	"github.com/netmoth/netmoth/internal/config"
)

// AFPacketStrategy  is...
type AFPacketStrategy struct {
	handles []*afpacket.TPacket
}

// New is ...
func (s *AFPacketStrategy) New(c *config.Config) ([]PacketDataSource, error) {
	var err error
	var compiledBpf []bpf.RawInstruction
	if len(c.Bpf) != 0 {
		compiledBpf, err = s.compileBpf(c.SnapLen, c.Bpf)
		if err != nil {
			return nil, err
		}
	}

	var res []PacketDataSource
	for i := 0; i < c.NumberOfRings; i++ {
		handle, err := afpacket.NewTPacket(
			afpacket.OptInterface(c.Interface),
			afpacket.OptFrameSize(c.SnapLen),
			afpacket.OptTPacketVersion(afpacket.TPacketVersion3),
		)
		if err != nil {
			return nil, err
		}
		if c.NumberOfRings > 1 {
			if err = handle.SetFanout(afpacket.FanoutHash, uint16(clusterID)); err != nil {
				return nil, err
			}
		}
		if len(c.Bpf) != 0 {
			if err = handle.SetBPF(compiledBpf); err != nil {
				return nil, err
			}
		}

		s.handles = append(s.handles, handle)
		res = append(res, handle)
	}

	return res, nil
}

// Destroy is ...
func (s *AFPacketStrategy) Destroy() {
	for _, handle := range s.handles {
		handle.Close()
	}
}

// PacketStats  is ...
func (s *AFPacketStrategy) PacketStats() (received uint64, dropped uint64) {
	for _, handle := range s.handles {
		_, stats, err := handle.SocketStats()
		if err != nil {
			continue
		}
		received += uint64(stats.Packets())
		dropped += uint64(stats.Drops())
	}
	return
}

func (s *AFPacketStrategy) compileBpf(snapLen int, bpfFilter string) ([]bpf.RawInstruction, error) {
	instructions, err := pcap.CompileBPFFilter(layers.LinkTypeEthernet, snapLen, bpfFilter)
	if err != nil {
		return nil, err
	}
	var res []bpf.RawInstruction
	for _, i := range instructions {
		res = append(res, bpf.RawInstruction{
			Op: i.Code,
			Jt: i.Jt,
			Jf: i.Jf,
			K:  i.K,
		})
	}
	return res, nil
}
