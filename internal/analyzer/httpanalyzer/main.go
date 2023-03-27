package httpanalyzer

import (
	"bufio"
	"fmt"
	"io"
	"net/http"

	"github.com/netmoth/netmoth/internal/connection"
	"github.com/netmoth/netmoth/internal/signature"
)

// Analyze is ...
func Analyze(conn *connection.Connection, detector signature.Detector) (*HTTP, error) {
	h := new(HTTP)
	reader := bufio.NewReader(conn.Payload)
	for {
		req, err := http.ReadRequest(reader)
		if err == io.EOF {
			 return nil, nil
		} else if err != nil {
			return nil, fmt.Errorf("error, ReadRequest, %w", err)
		}

		h.Request = Request{
			Method:           req.Method,
			URL:              req.URL.String(),
			Headers:          req.Header,
			ContentLength:    req.ContentLength,
			TransferEncoding: req.TransferEncoding,
			Host:             req.Host,
		}

		url := h.Request.URL
		if url == "/" {
			url = ""
		}

		h.Response, err = detector.Scan(&signature.Request{
			IP:         conn.DestinationIP,
			Port:       conn.DestinationPort,
			TrackerURL: h.Request.Host + url,
		})
		if err != nil {
			return nil, err
		}
		return h, nil
	}
}
