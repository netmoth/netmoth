package utils

import (
	"errors"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/google/gopacket"

	"github.com/google/gopacket/pcap"
)

// CheckIfInterfaceExists is ...
func CheckIfInterfaceExists(iface string) error {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	for _, device := range devices {
		if device.Name == iface {
			return nil
		}
	}
	return errors.New("specified network interface does not exist")
}

// InterfaceAddresses is ...
func InterfaceAddresses(interfaceName string) (addresses []string) {
	i, err := net.InterfaceByName(interfaceName)
	if err != nil {
		panic("interface invalid")
	}

	addrs, err := i.Addrs()
	if err != nil {
		panic("interface invalid")
	}

	for _, addr := range addrs {
		_addr := strings.Split(addr.String(), "/")
		addresses = append(addresses, _addr[0])
	}

	return addresses
}

// ProcessPorts is ...
func ProcessPorts(transport gopacket.Flow) (srcPort, dstPort int) {
	srcPort, _ = strconv.Atoi(transport.Src().String())
	dstPort, _ = strconv.Atoi(transport.Dst().String())
	return srcPort, dstPort
}
