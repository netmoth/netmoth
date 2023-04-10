package dnsanalyzer

// SOA is ...
type SOA struct {
	MName   string
	RName   string
	Serial  uint32
	Refresh uint32
	Retry   uint32
	Expire  uint32
	TTL     uint32
}

// Question is ...
type Question struct {
	Name  string
	Type  string
	Class string
}

// Record is ...
type Record struct {
	Name  string
	Type  string
	Class string
	Data  string   `json:",omitempty"`
	IP    string   `json:",omitempty"`
	NS    string   `json:",omitempty"`
	CNAME string   `json:",omitempty"`
	PTR   string   `json:",omitempty"`
	TXT   []string `json:",omitempty"`
	SOA   SOA      `json:",omitempty"`
	TTL   uint32
}

// DNS is ...
type DNS struct {
	OpCode       string
	ResponseCode string
	Questions    []Question
	Answers      []Record
	Authorities  []Record
	Additionals  []Record
	ID           uint16
	QR           bool
	AA           bool
	TC           bool
}

// Key is ...
func (d *DNS) Key() string {
	return "dns"
}
