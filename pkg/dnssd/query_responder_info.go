package dnssd

import (
	"github.com/galenliu/chip/pkg/dnssd/responder"
	"time"
)

type QueryResponderInfo struct {
	Responder                 responder.RecordResponder
	reportNowAsAdditional     bool   // report as additional data required
	alsoReportAdditionalQName bool   // report more data when this record is listed
	additionalQName           string // if alsoReportAdditionalQName is set, send this extra data
	reportService             bool
	lastMulticastTime         time.Time
}

func NewQueryResponderInfo(r responder.RecordResponder) *QueryResponderInfo {
	return &QueryResponderInfo{
		Responder: r,
	}
}
