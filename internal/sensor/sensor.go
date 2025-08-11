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
	// Agent mode fields
	agentClient    *AgentClient
	agentMode      bool
	dataInterval   time.Duration
	healthInterval time.Duration
	// Ring buffer for connections
	connRing []*connection.Connection
	connHead int
	connSize int
	// Ring buffer for signatures
	sigRing     []signature.Detect
	sigHead     int
	sigSize     int
	bufferMutex sync.RWMutex
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

	// Initialize agent mode if enabled
	if config.AgentMode {
		s.agentMode = true
		s.dataInterval = time.Duration(config.DataInterval) * time.Second
		s.healthInterval = time.Duration(config.HealthInterval) * time.Second

		if s.dataInterval == 0 {
			s.dataInterval = 60 * time.Second // default 1 minute
		}
		if s.healthInterval == 0 {
			s.healthInterval = 300 * time.Second // default 5 minutes
		}

		s.agentClient = NewAgentClient(config.ServerURL, config.AgentID, config.Interface, config.AgentToken)

		// Pre-allocate ring buffers to avoid reallocations
		const maxConnBuffer = 100000
		const maxSigBuffer = 50000
		s.connRing = make([]*connection.Connection, maxConnBuffer)
		s.sigRing = make([]signature.Detect, maxSigBuffer)

		// Register agent with server
		if err := s.agentClient.Register(config.Interface); err != nil {
			log.Printf("Warning: Failed to register agent: %v", err)
		}

		// Start agent goroutines
		go s.agentDataSender()
		go s.agentHealthChecker()
	} else {
		// Only connect to database if not in agent mode
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

		s.detector = signature.New(*s.db)
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

	nl := packet.NetworkLayer()
	tl := packet.TransportLayer()
	if tl != nil {
		switch tl.LayerType() {
		case layers.LayerTypeTCP:
			if nl != nil {
				s.streamFactory.NewPacket(nl.NetworkFlow(), tl.(*layers.TCP))
			}
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
		// Analyze the connection
		if err := s.analyze(conn); err != nil {
			// Log the error but continue processing
			// log.Printf("Error analyzing connection: %v", err)
		}

		// Save to log
		logger.save(*conn)

		// Add to agent buffer for sending to server
		s.addToBuffer(conn)

		// Return connection to pool
		connection.GlobalConnectionPool.Put(conn)
	}
}

// agentDataSender periodically sends data to the central server
func (s *sensor) agentDataSender() {
	ticker := time.NewTicker(s.dataInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Snapshot ring buffers with minimal copying
			s.bufferMutex.Lock()
			connections := make([]*connection.Connection, s.connSize)
			for i := 0; i < s.connSize; i++ {
				idx := (s.connHead + i) % len(s.connRing)
				connections[i] = s.connRing[idx]
			}
			signatures := make([]signature.Detect, s.sigSize)
			for i := 0; i < s.sigSize; i++ {
				idx := (s.sigHead + i) % len(s.sigRing)
				signatures[i] = s.sigRing[idx]
			}
			s.bufferMutex.Unlock()

			if len(connections) > 0 || len(signatures) > 0 {
				stats := AgentStats{
					PacketsReceived:  s.packetStats.received,
					PacketsDropped:   s.packetStats.dropped,
					PacketsProcessed: s.packetStats.processed,
					ConnectionsFound: uint64(len(connections)),
				}

				if err := s.agentClient.SendData(connections, signatures, stats, s.sensorMeta.NetworkInterface); err != nil {
					log.Printf("Failed to send data to server: %v", err)
				} else {
					// Clear sizes only on success (constant time)
					s.bufferMutex.Lock()
					s.connSize = 0
					s.sigSize = 0
					s.bufferMutex.Unlock()
				}
			}
		}
	}
}

// agentHealthChecker periodically sends health checks
func (s *sensor) agentHealthChecker() {
	ticker := time.NewTicker(s.healthInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.agentClient.SendHealth(); err != nil {
				log.Printf("Health check failed: %v", err)
			}
		}
	}
}

// addToBuffer adds connection to ring buffer for sending
func (s *sensor) addToBuffer(conn *connection.Connection) {
	if !s.agentMode {
		return
	}

	s.bufferMutex.Lock()
	if len(s.connRing) == 0 {
		// Agent mode not initialized yet
		s.bufferMutex.Unlock()
		return
	}
	if s.connSize < len(s.connRing) {
		idx := (s.connHead + s.connSize) % len(s.connRing)
		s.connRing[idx] = conn
		s.connSize++
	} else {
		// Ring full: drop oldest by advancing head, write at tail
		s.connHead = (s.connHead + 1) % len(s.connRing)
		tail := (s.connHead + s.connSize - 1) % len(s.connRing)
		s.connRing[tail] = conn
	}
	s.bufferMutex.Unlock()
}

// addSignatureToBuffer adds signature to ring buffer for sending
func (s *sensor) addSignatureToBuffer(sig signature.Detect) {
	if !s.agentMode {
		return
	}

	s.bufferMutex.Lock()
	if len(s.sigRing) == 0 {
		s.bufferMutex.Unlock()
		return
	}
	if s.sigSize < len(s.sigRing) {
		idx := (s.sigHead + s.sigSize) % len(s.sigRing)
		s.sigRing[idx] = sig
		s.sigSize++
	} else {
		// Ring full
		s.sigHead = (s.sigHead + 1) % len(s.sigRing)
		tail := (s.sigHead + s.sigSize - 1) % len(s.sigRing)
		s.sigRing[tail] = sig
	}
	s.bufferMutex.Unlock()
}
