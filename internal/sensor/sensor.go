package sensor

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"github.com/netmoth/netmoth/internal/config"
	"github.com/netmoth/netmoth/internal/connection"
	"github.com/netmoth/netmoth/internal/sensor/strategies"
	"github.com/netmoth/netmoth/internal/signature"
	"github.com/netmoth/netmoth/internal/storage/postgres"
	"github.com/netmoth/netmoth/internal/utils"
)

type sensor struct {
	db       *postgres.Connect
	detector signature.Detector

	strategy   strategies.PacketsCaptureStrategy
	packets    []strategies.PacketDataSource
	sensorMeta *Metadata

	streamFactory *connection.TCPStreamFactory
	connections   chan *connection.Connection
}

// Metadata is ...
type Metadata struct {
	NetworkInterface string
	NetworkAddress   []string
}

// New is the entry point for analyzer
func New(config *config.Config) {
	var err error
	s := new(sensor)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pgDSN := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		config.Postgres.User,
		config.Postgres.Password,
		config.Postgres.Host,
		config.Postgres.DB,
	)

	s.db, err = postgres.New(ctx, &postgres.PgSQLConfig{
		DSN:             pgDSN,
		MaxConn:         50,
		MaxIdleConn:     10,
		MaxLifetimeConn: 300,
	})
	if err != nil {
		log.Fatal(err)
	}

	s.sensorMeta = &Metadata{
		NetworkInterface: config.Interface,
		NetworkAddress:   utils.InterfaceAddresses(config.Interface),
	}

	logSave, err := newSaver(config.LogFile, s.sensorMeta)
	if err != nil {
		log.Fatal(err)
	}

	s.detector = signature.New(*s.db)

	// add signatures in database
	//go s.detector.Update()

	var ok bool
	s.strategy, ok = strategies.Strategies()[config.Strategy]
	if !ok {
		os.Exit(1)
	}

	conn := make(chan *connection.Connection)
	s.connections = conn
	s.streamFactory = &connection.TCPStreamFactory{
		Connections: conn,
		ConnTimeout: config.ConnTimeout,
	}

	// init Assembler
	s.streamFactory.CreateAssembler()
	s.streamFactory.Ticker = time.NewTicker(time.Second * 10)

	s.packets, err = s.strategy.New(config)
	if err != nil {
		log.Fatal(err)
	}
	defer s.strategy.Destroy()

	exitCh := make(chan bool)
	for _, source := range s.packets {
		if config.ZeroCopy {
			go s.capturePacketsZeroCopy(source, exitCh)
		} else {
			go s.capturePackets(source, exitCh)
		}
	}
	go s.printStats(exitCh)
	go s.processConnections(logSave)

	defer close(exitCh)

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh
	signal.Reset(syscall.SIGINT, syscall.SIGTERM)
	signal.Stop(signalCh)
	log.Println("got signal, cleanup and exit...")
}

func (s *sensor) capturePacketsZeroCopy(source gopacket.ZeroCopyPacketDataSource, exit <-chan bool) {
	for {
		select {
		case <-exit:
			return
		default:
			// data, ci, err := source.ZeroCopyReadPacketData()
			data, _, err := source.ZeroCopyReadPacketData()
			if err != nil {
				continue
			}
			packet := gopacket.NewPacket(data, layers.LinkTypeEthernet, gopacket.DecodeStreamsAsDatagrams)
			//m := packet.Metadata()
			//m.CaptureInfo = ci
			//m.Truncated = m.Truncated || ci.CaptureLength < ci.Length
			go s.processPacket(packet)
		}
	}
}

func (s *sensor) capturePackets(source gopacket.PacketDataSource, exit <-chan bool) {
	packetSource := gopacket.NewPacketSource(source, layers.LinkTypeEthernet)
	packetSource.DecodeOptions = gopacket.NoCopy

	for {
		select {
		case <-exit:
			return
		default:
			packet, err := packetSource.NextPacket()
			if err != nil {
				continue
			}
			go s.processPacket(packet)
		}
	}
}

func (s *sensor) processPacket(packet gopacket.Packet) {
	ci := packet.Metadata().CaptureInfo

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

	//if eth := packet.LinkLayer(); eth != nil {
	//	srcMac := eth.LinkFlow().Src()
	//	fmt.Print(srcMac)
	//}

	//if ip := packet.NetworkLayer(); ip != nil {
	//	srcIp, dstIp := ip.NetworkFlow().Endpoints()
	//	fmt.Print(srcIp, dstIp)
	//}

	//if trans := packet.TransportLayer(); trans != nil {
	//	srcPort, dstPort := trans.TransportFlow().Endpoints()
	//	fmt.Print(srcPort, dstPort)
	//}
}

func (s *sensor) printStats(exit <-chan bool) {
	var receivedBefore uint64 = 0
	for {
		select {
		case <-exit:
			return
		default:
			received, dropped := s.strategy.PacketStats()
			pps := received - receivedBefore
			var packetLoss float64
			if received > 0 || dropped > 0 {
				packetLoss = float64(dropped) / float64(received+dropped) * 100
			}
			log.Printf("pps: %d, packet loss: %2f%%, goroutine number: %d\n", pps, packetLoss, runtime.NumGoroutine())
			receivedBefore = received
			time.Sleep(time.Second)
		}
	}
}

func (s *sensor) processConnections(logger *logSave) {
	for conn := range s.connections {
		if err := s.analyze(conn); err != nil {
			log.Println(err)
		}
		logger.save(*conn)
	}
}
