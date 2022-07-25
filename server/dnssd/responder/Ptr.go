package responder

import (
	"github.com/miekg/dns"
)

type PtrResponder struct {
	dns.PTR
}

func NewPtrResponder(name, target string) *PtrResponder {
	ptr := &PtrResponder{}
	ptr.Hdr.Name = name
	ptr.PTR.Ptr = target
	ptr.PTR.Hdr.Class = dns.ClassINET
	ptr.PTR.Hdr.Rrtype = dns.TypeAAAA
	return ptr
}

func (p *PtrResponder) GetClass() uint16 {
	return p.Hdr.Class
}

func (p *PtrResponder) GetName() string {
	return p.Hdr.Name
}

func (p *PtrResponder) GetType() uint16 {
	return p.Hdr.Rrtype
}

func (p *PtrResponder) GetTtl() uint32 {
	return p.Hdr.Ttl
}

func (p *PtrResponder) SetTtl(u uint32) {
	p.Hdr.Ttl = u
}

func (p *PtrResponder) GetRecord() *Record {
	return &Record{&p.PTR}
}
