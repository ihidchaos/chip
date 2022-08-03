package responder

import (
	"github.com/miekg/dns"
)

type IPv4Responder struct {
	dns.A
}

func NewIPv4Responder(qname string) *IPv4Responder {
	ip := &IPv4Responder{}
	ip.A = dns.A{
		Hdr: dns.RR_Header{
			Name:   qname,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
		},
	}
	ip.A.Hdr.Rrtype = dns.TypeA
	return ip
}

func (ip *IPv4Responder) GetTtl() uint32 {
	return ip.Hdr.Ttl
}

func (ip *IPv4Responder) SetTtl(ttl uint32) {
	ip.Header().Ttl = ttl
}

func (ip *IPv4Responder) GetQName() string {
	return ip.A.Hdr.Name
}

func (ip *IPv4Responder) GetQClass() uint16 {
	return ip.A.Hdr.Class
}

func (ip *IPv4Responder) GetQType() uint16 {
	return ip.A.Hdr.Rrtype
}

func (ip *IPv4Responder) GetRecord() *Record {
	return NewRecord(&ip.A)
}

type AddResponder interface {
	Append(rr dns.RR)
}

func (ip *IPv4Responder) AddRecord(r AddResponder) {
	r.Append(ip.GetRecord())
}
