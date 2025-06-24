package sensor

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
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
	strategy      strategies.PacketsCaptureStrategy
	db            *postgres.Connect
	detector      signature.Detector
	sensorMeta    *Metadata
	streamFactory *connection.TCPStreamFactory
	connections   chan *connection.Connection
	packets       []strategies.PacketDataSource
	packetPool    sync.Pool
	workerPool    chan struct{}
	statsMutex    sync.RWMutex
	packetStats   struct {
		received  uint64
		dropped   uint64
		processed uint64
	}
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

	s.packetPool = sync.Pool{
		New: func() any {
			return make([]byte, 0, config.SnapLen)
		},
	}

	workerCount := runtime.NumCPU() * 2
	if config.MaxCores > 0 {
		workerCount = config.MaxCores * 2
	}
	s.workerPool = make(chan struct{}, workerCount)

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

	netAddress, _ := utils.InterfaceAddresses(config.Interface)
	s.sensorMeta = &Metadata{
		NetworkInterface: config.Interface,
		NetworkAddress:   netAddress,
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

	conn := make(chan *connection.Connection, 10000)
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
			data, ci, err := source.ZeroCopyReadPacketData()
			if err != nil {
				s.statsMutex.Lock()
				s.packetStats.dropped++
				s.statsMutex.Unlock()
				continue
			}

			s.statsMutex.Lock()
			s.packetStats.received++
			s.statsMutex.Unlock()

			packet := gopacket.NewPacket(data, layers.LinkTypeEthernet, gopacket.DecodeStreamsAsDatagrams)

			m := packet.Metadata()
			m.CaptureInfo = ci
			m.Truncated = m.Truncated || ci.CaptureLength < ci.Length

			select {
			case s.workerPool <- struct{}{}:
				go func() {
					defer func() { <-s.workerPool }()
					s.processPacket(packet)
					s.statsMutex.Lock()
					s.packetStats.processed++
					s.statsMutex.Unlock()
				}()
			default:
				s.processPacket(packet)
				s.statsMutex.Lock()
				s.packetStats.processed++
				s.statsMutex.Unlock()
			}
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
				s.statsMutex.Lock()
				s.packetStats.dropped++
				s.statsMutex.Unlock()
				continue
			}

			s.statsMutex.Lock()
			s.packetStats.received++
			s.statsMutex.Unlock()

			select {
			case s.workerPool <- struct{}{}:
				go func() {
					defer func() { <-s.workerPool }()
					s.processPacket(packet)
					s.statsMutex.Lock()
					s.packetStats.processed++
					s.statsMutex.Unlock()
				}()
			default:
				s.processPacket(packet)
				s.statsMutex.Lock()
				s.packetStats.processed++
				s.statsMutex.Unlock()
			}
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
			select {
			case s.connections <- udp:
			default:
				log.Printf("Connection channel full, dropping UDP packet")
			}
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
	var processedBefore uint64 = 0
	for {
		select {
		case <-exit:
			return
		case <-time.After(time.Second * 5):
			s.statsMutex.RLock()
			received := s.packetStats.received
			dropped := s.packetStats.dropped
			processed := s.packetStats.processed
			s.statsMutex.RUnlock()

			receivedDiff := received - receivedBefore
			processedDiff := processed - processedBefore
			receivedBefore = received
			processedBefore = processed

			log.Printf("Stats: Received: %d/s, Processed: %d/s, Total Received: %d, Total Dropped: %d, Total Processed: %d",
				receivedDiff/5, processedDiff/5, received, dropped, processed)
		}
	}
}

func (s *sensor) processConnections(logger *logSave) {
	for conn := range s.connections {
		// Анализируем соединение
		if err := s.analyze(conn); err != nil {
			// Логируем ошибку, но продолжаем обработку
			// log.Printf("Error analyzing connection: %v", err)
		}

		// Сохраняем в лог
		logger.save(*conn)

		// Возвращаем соединение в пул
		connection.GlobalConnectionPool.Put(conn)
	}
}
