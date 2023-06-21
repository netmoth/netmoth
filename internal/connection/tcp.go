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
		Analyzers:       make(map[string]interface{}),
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
}

// New is ...
func (tsf *TCPStreamFactory) New(n, t gopacket.Flow, tcp *layers.TCP, ac reassembly.AssemblerContext) reassembly.Stream {
	ts := &tcpStream{
		net:       n,
		transport: t,
		payload:   new(bytes.Buffer),
		startTime: ac.GetCaptureInfo().Timestamp,
		tcpState:  reassembly.NewTCPSimpleFSM(reassembly.TCPSimpleFSMOptions{}),
		done:      make(chan bool),
	}
	go func() {
		// wait for reassembly to be done
		<-ts.done
		// ignore empty streams
		if ts.packets > 0 {
			c := NewTCP(ts)
			tsf.Connections <- c
		}
	}()
	return ts
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
}
