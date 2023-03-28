package sensor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"github.com/netmoth/netmoth/internal/config"
	"github.com/netmoth/netmoth/internal/connection"
	"github.com/netmoth/netmoth/internal/signature"
	"github.com/netmoth/netmoth/internal/storage/postgres"
	"github.com/netmoth/netmoth/internal/utils"
)

var (
	detector   signature.Detector
	sensorMeta *Metadata
)

// Metadata is ...
type Metadata struct {
	NetworkInterface string
	NetworkAddress   []string
}

func getSensorMetadata(interfaceName string) *Metadata {
	return &Metadata{
		NetworkInterface: interfaceName,
		NetworkAddress:   utils.InterfaceAddresses(interfaceName),
	}
}

type sensor struct {
	source        gopacket.ZeroCopyPacketDataSource
	streamFactory *connection.TCPStreamFactory
	connections   chan *connection.Connection
}

// New is the entry point for analyzer
func New(ctx context.Context, config *config.Config) {
	sensorMeta = getSensorMetadata(config.Interface)

	if err := newSaver(config.LogFile); err != nil {
		log.Fatal(err)
	}

	pgDSN := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		config.Postgres.User,
		config.Postgres.Password,
		config.Postgres.Host,
		config.Postgres.DB,
	)

	db, err := postgres.New(ctx, &postgres.PgSQLConfig{
		DSN:             pgDSN,
		MaxConn:         50,
		MaxIdleConn:     10,
		MaxLifetimeConn: 300,
	})
	if err != nil {
		log.Fatal(err)
	}

	detector = signature.New(*db)

	// add signatures in database
	//go detector.Update()

	source, err := newLibpcap(config)
	if err != nil {
		log.Fatal(err)
	}

	conn := make(chan *connection.Connection)
	s := &sensor{
		source:      source,
		connections: conn,
		streamFactory: &connection.TCPStreamFactory{
			Connections: conn,
			ConnTimeout: config.ConnTimeout,
		},
	}

	// go s.processConnections()

	fmt.Printf("analyzer is running and logging to %s. Press CTL+C to stop...", logger.fileName)
	fmt.Println()

	s.run()
}

func (s *sensor) run() {
	s.streamFactory.CreateAssembler()
	s.streamFactory.Ticker = time.NewTicker(time.Second * 10)
	for {
		p, ci, err := s.source.ZeroCopyReadPacketData()
		if err != nil {
			log.Printf("s.run() return: %s", err)
			continue
		}
		packet := gopacket.NewPacket(p, layers.LayerTypeEthernet, gopacket.DecodeStreamsAsDatagrams)
		go s.processNewPacket(packet, ci)
	}
}

func (s *sensor) processNewPacket(packet gopacket.Packet, ci gopacket.CaptureInfo) {
	if packet.TransportLayer() != nil {
		layer := packet.TransportLayer()
		switch layer.LayerType() {
		case layers.LayerTypeTCP:
			s.streamFactory.NewPacket(packet.NetworkLayer().NetworkFlow(), packet.TransportLayer().(*layers.TCP))
			return
		case layers.LayerTypeUDP:
			udp := connection.NewUDP(packet, ci)
			s.connections <- udp
			return
		}
	}
}

func (s *sensor) processConnections() {
	for conn := range s.connections {
		if err := analyze(conn); err != nil {
			log.Println(err)
		}
		logger.save(*conn)
	}
}
