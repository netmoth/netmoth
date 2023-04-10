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
	Headers          http.Header
	Method           string
	URL              string
	Host             string
	TransferEncoding []string
	ContentLength    int64
}

// Key is ..
func (h *HTTP) Key() string {
	return "http"
}
