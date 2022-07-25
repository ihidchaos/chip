package responder

import (
	"github.com/miekg/dns"
)

type IPv6Responder struct {
	dns.AAAA
}

func NewIPv6Responder(qname string) *IPv6Responder {
	ip6 := &IPv6Responder{}
	ip6.Hdr.Name = qname
	ip6.Hdr.Class = dns.ClassINET
	ip6.Hdr.Rrtype = dns.TypeAAAA
	return ip6
}

func (ipv6 IPv6Responder) GetRecord() *Record {
	return NewRecord(&ipv6.AAAA)
}

func (ipv6 IPv6Responder) GetClass() uint16 {
	return ipv6.Header().Class
}

func (ipv6 IPv6Responder) GetName() string {
	return ipv6.Header().Name
}

func (ipv6 IPv6Responder) GetType() uint16 {
	return ipv6.Hdr.Rrtype
}

func (ipv6 IPv6Responder) GetTtl() uint32 {
	return ipv6.Hdr.Ttl
}

func (ipv6 IPv6Responder) SetTtl(u uint32) {
	ipv6.Hdr.Ttl = u
}
