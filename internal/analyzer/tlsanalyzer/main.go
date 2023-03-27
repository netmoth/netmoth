package tlsanalyzer

import (
	"github.com/netmoth/netmoth/internal/connection"
)

// Analyze is ...
func Analyze(c *connection.Connection) (*TLS, error) {
	result := TLS{
		Str: c.Payload.Len(),
	}
	return &result, nil
}
