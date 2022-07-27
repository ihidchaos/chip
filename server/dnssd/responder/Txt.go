package responder

import (
	"github.com/miekg/dns"
)

const kTxtDefaultTtl = 4500

func NewTxtResponder(qname string, txt []string) *TxtResponder {
	t := &TxtResponder{}
	t.TXT = dns.TXT{
		Hdr: dns.RR_Header{
			Name:   qname,
			Rrtype: dns.TypeTXT,
			Class:  dns.ClassINET,
		},
		Txt: txt,
	}
	return t
}

type TxtResponder struct {
	dns.TXT
}

func (t *TxtResponder) GetRecord() *Record {
	return NewRecord(&t.TXT)
}

func (t *TxtResponder) GetQClass() uint16 {
	return t.Hdr.Class
}

func (t *TxtResponder) GetQName() string {
	return t.Hdr.Name
}

func (t *TxtResponder) GetQType() uint16 {
	return t.Hdr.Rrtype
}

func (t *TxtResponder) GetTtl() uint32 {
	return t.Hdr.Ttl
}

func (t *TxtResponder) SetTtl(u uint32) {
	t.Hdr.Ttl = u
}
