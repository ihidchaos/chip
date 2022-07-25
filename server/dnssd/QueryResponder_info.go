package dnssd

import (
	"github.com/galenliu/chip/server/dnssd/responder"
	"time"
)

type QueryResponderInfo struct {
	Responder                 responder.RecordResponder
	reportNowAsAdditional     bool
	alsoReportAdditionalQName bool
	additionalQName           string
	reportService             bool
	LastMulticastTime         time.Time
}

func NewQueryResponderInfo(r responder.RecordResponder) QueryResponderInfo {
	return QueryResponderInfo{
		Responder: r,
	}
}
