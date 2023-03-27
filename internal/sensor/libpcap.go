package sensor

import (
	"github.com/google/gopacket/pcap"

	"github.com/netmoth/netmoth/internal/config"
)

func newLibpcap(c *config.Config) (*pcap.Handle, error) {
	handle, err := pcap.OpenLive(c.Interface, int32(c.SnapLen), c.Promiscuous, pcap.BlockForever)
	if err != nil {
		return nil, err
	}
	err = handle.SetBPFFilter(c.Bpf)
	if err != nil {
		return nil, err
	}
	return handle, nil
}
