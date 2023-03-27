package dnsanalyzer

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"github.com/netmoth/netmoth/internal/connection"
)

// Analyze is ...
func Analyze(c *connection.Connection) (*DNS, error) {
	d := new(DNS)
	var dns layers.DNS
	var decoded []gopacket.LayerType

	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeDNS, &dns)
	if err := parser.DecodeLayers(c.Payload.Bytes(), &decoded); err != nil {
		return nil, err
	}
	for _, layerType := range decoded {
		if layerType == layers.LayerTypeDNS {
			d = newDNSResult(dns)
		}
	}
	return d, nil
}
