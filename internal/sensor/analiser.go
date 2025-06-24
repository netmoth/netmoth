package sensor

import (
	"github.com/netmoth/netmoth/internal/analyzer/contentanalyzer"
	"github.com/netmoth/netmoth/internal/analyzer/dnsanalyzer"
	"github.com/netmoth/netmoth/internal/analyzer/http2analyzer"
	"github.com/netmoth/netmoth/internal/analyzer/httpanalyzer"
	"github.com/netmoth/netmoth/internal/analyzer/tlsanalyzer"
	"github.com/netmoth/netmoth/internal/connection"
)

func (s *sensor) analyze(conn *connection.Connection) error {
	// Always try content analysis first
	contentResult, err := contentanalyzer.Analyze(conn)
	if err == nil && contentResult != nil {
		conn.Analyzers[contentResult.Key()] = contentResult
	}

	switch {
	case conn.SourcePort == 80 || conn.DestinationPort == 80:
		result, err := httpanalyzer.Analyze(conn, s.detector)
		if err != nil {
			return err
		}
		conn.Analyzers[result.Key()] = result

	case conn.SourcePort == 443 || conn.DestinationPort == 443 || conn.SourcePort == 8443 || conn.DestinationPort == 8443:
		// First try TLS analysis
		tlsResult, err := tlsanalyzer.Analyze(conn)
		if err == nil && tlsResult != nil {
			conn.Analyzers[tlsResult.Key()] = tlsResult

			// If TLS analysis detected HTTP/2, try HTTP/2 analysis
			if tlsResult.Extensions["detected_protocol"] == "HTTP/2" {
				http2Result, err := http2analyzer.Analyze(conn)
				if err == nil && http2Result != nil {
					conn.Analyzers[http2Result.Key()] = http2Result
				}
			}
		} else {
			// If TLS analysis failed, try HTTP/2 analysis directly
			http2Result, err := http2analyzer.Analyze(conn)
			if err == nil && http2Result != nil {
				conn.Analyzers[http2Result.Key()] = http2Result
			}
		}

	case conn.SourcePort == 53 || conn.DestinationPort == 53:
		result, err := dnsanalyzer.Analyze(conn)
		if err != nil {
			return err
		}
		conn.Analyzers[result.Key()] = result
	}

	return nil
}
