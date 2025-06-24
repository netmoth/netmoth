package connection

import (
	"bytes"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/reassembly"

	"github.com/netmoth/netmoth/internal/utils"
)

type tcpStream struct {
	startTime time.Time
	payload   *bytes.Buffer
	tcpState  *reassembly.TCPSimpleFSM
	done      chan bool
	net       gopacket.Flow
	transport gopacket.Flow
	duration  time.Duration
	packets   int
	reused    bool
}

// NewTCP is ...
func NewTCP(ts *tcpStream) *Connection {
	srcPort, dstPort, _ := utils.ProcessPorts(ts.transport)
	return &Connection{
		Timestamp:       ts.startTime,
		UID:             ts.net.FastHash() + ts.transport.FastHash(),
		SourceIP:        ts.net.Src().String(),
		SourcePort:      srcPort,
		DestinationIP:   ts.net.Dst().String(),
		DestinationPort: dstPort,
		TransportType:   "tcp",
		Duration:        ts.duration.Seconds(),
		State:           ts.tcpState.String(),
		Payload:         ts.payload,
		Analyzers:       make(map[string]any),
	}
}

// Accept is ...
func (ts *tcpStream) Accept(tcp *layers.TCP, ci gopacket.CaptureInfo, dir reassembly.TCPFlowDirection, nextSeq reassembly.Sequence, start *bool, ac reassembly.AssemblerContext) bool {
	tempDuration := ci.Timestamp.Sub(ts.startTime)
	if tempDuration.Seconds() > ts.duration.Seconds() {
		ts.duration = tempDuration
	}
	ts.tcpState.CheckState(tcp, dir)
	return true
}

// ReassembledSG is ...
func (ts *tcpStream) ReassembledSG(sg reassembly.ScatterGather, ac reassembly.AssemblerContext) {
	length, _ := sg.Lengths()
	data := sg.Fetch(length)
	if length > 0 {
		ts.payload.Write(data)
	}
	ts.packets++
}

// ReassemblyComplete is ...
func (ts *tcpStream) ReassemblyComplete(ac reassembly.AssemblerContext) bool {
	ts.done <- true
	return false
}

// TCPStreamFactory is ...
type TCPStreamFactory struct {
	Assembler      *reassembly.Assembler
	Ticker         *time.Ticker
	Connections    chan *Connection
	ConnTimeout    int
	AssemblerMutex sync.Mutex
	streamPool     sync.Pool
	bufferPool     sync.Pool
}

// New is ...
func (tsf *TCPStreamFactory) New(n, t gopacket.Flow, tcp *layers.TCP, ac reassembly.AssemblerContext) reassembly.Stream {
	tsInterface := tsf.streamPool.Get()
	var ts *tcpStream
	if tsInterface == nil {
		ts = &tcpStream{
			payload:  tsf.getBuffer(),
			tcpState: reassembly.NewTCPSimpleFSM(reassembly.TCPSimpleFSMOptions{}),
			done:     make(chan bool, 1),
		}
	} else {
		ts = tsInterface.(*tcpStream)
		ts.reused = true
		ts.payload.Reset()
		ts.tcpState = reassembly.NewTCPSimpleFSM(reassembly.TCPSimpleFSMOptions{})
		ts.packets = 0
		ts.duration = 0
	}

	ts.net = n
	ts.transport = t
	ts.startTime = ac.GetCaptureInfo().Timestamp

	go func() {
		<-ts.done
		if ts.packets > 0 {
			c := NewTCP(ts)
			select {
			case tsf.Connections <- c:
			default:
				// Если канал переполнен, логируем и пропускаем
				// log.Printf("Connection channel full, dropping TCP connection")
			}
		}
		tsf.returnStream(ts)
	}()
	return ts
}

// getBuffer получает буфер из пула или создает новый
func (tsf *TCPStreamFactory) getBuffer() *bytes.Buffer {
	bufferInterface := tsf.bufferPool.Get()
	if bufferInterface == nil {
		return new(bytes.Buffer)
	}
	buffer := bufferInterface.(*bytes.Buffer)
	buffer.Reset()
	return buffer
}

// returnStream возвращает stream в пул
func (tsf *TCPStreamFactory) returnStream(ts *tcpStream) {
	if ts.reused {
		tsf.bufferPool.Put(ts.payload)
		tsf.streamPool.Put(ts)
	}
}

// NewPacket is ...
func (tsf *TCPStreamFactory) NewPacket(netFlow gopacket.Flow, tcp *layers.TCP) {
	select {
	case <-tsf.Ticker.C:
		tsf.AssemblerMutex.Lock()
		tsf.Assembler.FlushCloseOlderThan(time.Now().Add(time.Second * time.Duration(-1*tsf.ConnTimeout)))
		tsf.AssemblerMutex.Unlock()
	default:
		// pass through
	}
	tsf.AssemblePacket(netFlow, tcp)
}

// AssemblePacket is ...
func (tsf *TCPStreamFactory) AssemblePacket(netFlow gopacket.Flow, tcp *layers.TCP) {
	tsf.AssemblerMutex.Lock()
	tsf.Assembler.Assemble(netFlow, tcp)
	tsf.AssemblerMutex.Unlock()
}

// CreateAssembler is ...
func (tsf *TCPStreamFactory) CreateAssembler() {
	streamPool := reassembly.NewStreamPool(tsf)
	tsf.Assembler = reassembly.NewAssembler(streamPool)
	tsf.Assembler.MaxBufferedPagesTotal = 100000
	tsf.Assembler.MaxBufferedPagesPerConnection = 1000
}
