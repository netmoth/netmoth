package dnsanalyzer

import "github.com/google/gopacket/layers"

func newDNSResult(parsedDNS layers.DNS) *DNS {
	return &DNS{
		ID:           parsedDNS.ID,
		QR:           parsedDNS.QR,
		OpCode:       parsedDNS.OpCode.String(),
		AA:           parsedDNS.AA,
		TC:           parsedDNS.TC,
		ResponseCode: parsedDNS.ResponseCode.String(),
		Questions:    newDNSQuestions(parsedDNS),
		Answers:      newDNSRecords(parsedDNS.Answers),
		Authorities:  newDNSRecords(parsedDNS.Authorities),
		Additionals:  newDNSRecords(parsedDNS.Additionals),
	}
}

func newDNSQuestions(parsedDNS layers.DNS) []Question {
	var dnsQuestions []Question

	for _, question := range parsedDNS.Questions {
		dnsQuestions = append(dnsQuestions, newDNSQuestion(question))
	}
	return dnsQuestions
}

func newDNSQuestion(question layers.DNSQuestion) Question {
	return Question{
		Name:  string(question.Name),
		Class: question.Class.String(),
		Type:  question.Type.String(),
	}
}

func newDNSRecords(records []layers.DNSResourceRecord) []Record {
	var dnsRecords []Record

	for _, record := range records {
		dnsRecords = append(dnsRecords, newDNSRecord(record))
	}
	return dnsRecords
}

func newDNSRecord(record layers.DNSResourceRecord) Record {
	return Record{
		Name:  string(record.Name),
		Type:  record.Type.String(),
		Class: record.Class.String(),
		Data:  string(record.Data),
		IP:    record.IP.String(),
		NS:    string(record.NS),
		CNAME: string(record.CNAME),
		PTR:   string(record.PTR),
		TXT:   convertDNSTXTToStrings(record.TXTs),
		SOA:   newDNSSOA(record.SOA),
	}
}

func convertDNSTXTToStrings(txtBytes [][]byte) []string {
	var txtStrings []string

	for _, txt := range txtBytes {
		txtStrings = append(txtStrings, string(txt))
	}
	return txtStrings
}

func newDNSSOA(soa layers.DNSSOA) SOA {
	return SOA{
		MName:   string(soa.MName),
		RName:   string(soa.RName),
		Expire:  soa.Expire,
		Refresh: soa.Refresh,
		Serial:  soa.Serial,
		Retry:   soa.Retry,
		TTL:     soa.Minimum,
	}
}
