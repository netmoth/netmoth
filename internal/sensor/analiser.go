package sensor

import (
	"github.com/netmoth/netmoth/internal/analyzer/dnsanalyzer"
	"github.com/netmoth/netmoth/internal/analyzer/httpanalyzer"
	"github.com/netmoth/netmoth/internal/analyzer/tlsanalyzer"
	"github.com/netmoth/netmoth/internal/connection"
)

func (s *sensor) analyze(conn *connection.Connection) error {
	switch {
	case conn.SourcePort == 80 || conn.DestinationPort == 80:
		result, err := httpanalyzer.Analyze(conn, s.detector)
		if err != nil {
			return err
		}
		conn.Analyzers[result.Key()] = result

	case conn.SourcePort == 443 || conn.DestinationPort == 443:
		result, err := tlsanalyzer.Analyze(conn)
		if err != nil {
			return err
		}
		conn.Analyzers[result.Key()] = result

	case conn.SourcePort == 53 || conn.DestinationPort == 53:
		result, err := dnsanalyzer.Analyze(conn)
		if err != nil {
			return err
		}
		conn.Analyzers[result.Key()] = result
	}

	return nil
}
