package responder

import (
	"github.com/miekg/dns"
)

type SrvResponder struct {
	dns.SRV
}

func NewSrvResponder(qName string, serverName string, port uint16) *SrvResponder {
	ser := &SrvResponder{}
	ser.SRV = dns.SRV{
		Hdr: dns.RR_Header{
			Name:   qName,
			Rrtype: dns.TypeSRV,
			Class:  dns.ClassINET,
		},
		Port:   port,
		Target: serverName,
	}
	return ser
}

func (s *SrvResponder) GetRecord() *Record {
	return NewRecord(&s.SRV)
}

func (s SrvResponder) GetClass() uint16 {
	return s.Hdr.Class
}

func (s SrvResponder) GetName() string {
	return s.Hdr.Name
}

func (s SrvResponder) GetType() uint16 {
	return s.Hdr.Rrtype
}

func (s SrvResponder) GetTtl() uint32 {
	return s.Hdr.Ttl
}

func (s *SrvResponder) SetTtl(u uint32) {
	s.Hdr.Ttl = u
}
