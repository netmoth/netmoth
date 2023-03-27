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
	TTL   uint32
	Data  string   `json:",omitempty"`
	IP    string   `json:",omitempty"`
	NS    string   `json:",omitempty"`
	CNAME string   `json:",omitempty"`
	PTR   string   `json:",omitempty"`
	TXT   []string `json:",omitempty"`
	SOA   SOA      `json:",omitempty"`
}

// DNS is ...
type DNS struct {
	ID           uint16
	QR           bool
	OpCode       string
	AA           bool
	TC           bool
	ResponseCode string
	Questions    []Question
	Answers      []Record
	Authorities  []Record
	Additionals  []Record
}

// Key is ...
func (d *DNS) Key() string {
	return "dns"
}
