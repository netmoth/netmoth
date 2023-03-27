package httpanalyzer

import (
	"net/http"

	"github.com/netmoth/netmoth/internal/signature"
)

// HTTP is ...
type HTTP struct {
	Request  Request
	Response []signature.Detect
}

// Request is ...
type Request struct {
	Method           string
	URL              string
	Headers          http.Header
	ContentLength    int64
	TransferEncoding []string
	Host             string
}

// Key is ..
func (h *HTTP) Key() string {
	return "http"
}
