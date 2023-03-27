package tlsanalyzer

// TLS is ...
type TLS struct {
	Str int
}

// Key is ...
func (sr *TLS) Key() string {
	return "tls"
}
