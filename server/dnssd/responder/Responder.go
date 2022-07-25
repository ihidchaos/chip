package responder

import (
	"github.com/miekg/dns"
)

const DefaultTtl uint32 = 120

type RecordResponder interface {
	GetClass() uint16
	GetName() string
	GetType() uint16
	GetTtl() uint32
	SetTtl(uint32)
	RecordProvider
}

type RecordProvider interface {
	GetRecord() *Record
}

type Record struct {
	dns.RR
}

func NewRecord(rr dns.RR) *Record {
	return &Record{
		RR: rr,
	}
}

func (rr *Record) SetTtl(ttl uint32) {
	rr.Header().Ttl = ttl
}
