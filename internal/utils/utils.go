package utils

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// CheckIfInterfaceExists checks if the specified network interface exists.
// It takes the name of the interface as a string and returns an error if the interface does not exist.
func CheckIfInterfaceExists(iface string) error {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	deviceMap := make(map[string]bool)
	for _, device := range devices {
		deviceMap[device.Name] = true
	}
	if deviceMap[iface] {
		return nil
	}
	return errors.New("specified network interface does not exist")
}

// InterfaceAddresses returns the IP addresses associated with a network interface.
// It takes an interface name as input and returns a slice of strings containing the IP addresses,
// or an error if there was a problem retrieving the addresses.
func InterfaceAddresses(interfaceName string) ([]string, error) {
	i, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return nil, fmt.Errorf("invalid interface: %v", err)
	}

	addrs, err := i.Addrs()
	if err != nil {
		return nil, fmt.Errorf("failed to get addresses: %v", err)
	}

	var addresses []string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			addresses = append(addresses, ipnet.IP.String())
		}
	}

	return addresses, nil
}

// ProcessPorts takes a gopacket.Flow as input and returns the source port and destination port as integers.
// If there is an error in converting the string to integer, it will return an error.
func ProcessPorts(transport gopacket.Flow) (srcPort, dstPort int, err error) {
	srcPort, err = strconv.Atoi(transport.Src().String())
	if err != nil {
		return 0, 0, err
	}

	dstPort, err = strconv.Atoi(transport.Dst().String())
	if err != nil {
		return 0, 0, err
	}

	return srcPort, dstPort, nil
}
