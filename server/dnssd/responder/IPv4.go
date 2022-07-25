package responder

import (
	"github.com/miekg/dns"
	"net"
	"net/netip"
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

func (ip *IPv4Responder) AddAllResponses(srcAddrPort, destAddrPort netip.AddrPort, interfaceId net.Interface, delegate Delegate, configuration *ResponseConfiguration) {
	//TODO implement me
	panic("implement me")
}

func (ip *IPv4Responder) SetTtl(ttl uint32) {
	ip.Header().Ttl = ttl
}

func (ip *IPv4Responder) GetName() string {
	return ip.A.Hdr.Name
}

func (ip *IPv4Responder) GetClass() uint16 {
	return ip.A.Hdr.Class
}

func (ip *IPv4Responder) GetType() uint16 {
	return ip.A.Hdr.Rrtype
}

func (ip *IPv4Responder) GetRecord() *Record {
	return &Record{&ip.A}
}
