package httpanalyzer

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/netmoth/netmoth/internal/connection"
	"github.com/netmoth/netmoth/internal/signature"
)

func TestAnalyze_SimpleRequest(t *testing.T) {
	// HTTP/1.0 does not require Host header; avoids filling TrackerURL and DB calls
	raw := "GET / HTTP/1.0\r\n\r\n"
	conn := &connection.Connection{
		DestinationIP:   "",
		DestinationPort: 0,
		Payload:         bytes.NewBufferString(raw),
	}
	det := &signature.Detector{}
	res, err := Analyze(conn, *det)
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}
	if res == nil || res.Request.Method != http.MethodGet {
		t.Fatalf("unexpected result: %+v", res)
	}
}
