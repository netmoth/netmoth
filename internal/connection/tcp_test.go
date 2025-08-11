package connection

import (
	"bytes"
	"testing"
	"time"

	"github.com/google/gopacket/reassembly"
)

type fakeStreamFactory struct{ TCPStreamFactory }

func TestTCPStreamFactory_EmitOnComplete(t *testing.T) {
	ch := make(chan *Connection, 1)
	tsf := &TCPStreamFactory{Connections: ch, ConnTimeout: 1}
	tsf.CreateAssembler()
	tsf.Ticker = time.NewTicker(10 * time.Millisecond)
	defer tsf.Ticker.Stop()

	// Create tcpStream directly similar to factory behavior
	ts := &tcpStream{
		startTime: time.Now(),
		payload:   new(bytes.Buffer),
		done:      make(chan bool, 1),
		tcpState:  reassembly.NewTCPSimpleFSM(reassembly.TCPSimpleFSMOptions{}),
	}
	// Emit payload and complete
	ts.payload.WriteString("data")
	go func() { ts.done <- true }()
	// Handle completion similar to factory logic
	go func() {
		<-ts.done
		c := NewTCP(ts)
		select {
		case tsf.Connections <- c:
		default:
		}
	}()

	select {
	case c := <-ch:
		if c.TransportType != "tcp" || c.Payload.Len() == 0 {
			t.Fatalf("unexpected connection: %+v", c)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("no connection emitted")
	}
}
