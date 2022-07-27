package dnssd

import (
	"github.com/galenliu/chip/server/dnssd/responder"
	"time"
)

type QueryResponderRecordFilter struct {
	mIncludeAdditionalRepliesOnly bool
	mReplyFilter                  responder.ReplyFilter
	mIncludeOnlyMulticastBefore   time.Time
}

func NewQueryResponderRecordFilter() *QueryResponderRecordFilter {
	return &QueryResponderRecordFilter{}
}

func (f *QueryResponderRecordFilter) SetReplyFilter(filter responder.ReplyFilter) *QueryResponderRecordFilter {
	f.mReplyFilter = filter
	return f
}

func (f *QueryResponderRecordFilter) SetIncludeOnlyMulticastBeforeMS(t time.Time) {
	f.mIncludeOnlyMulticastBefore = t
}

func (f *QueryResponderRecordFilter) Accept(record *QueryResponderInfo) bool {
	if record.Responder == nil {
		return false
	}
	if f.mIncludeAdditionalRepliesOnly && !record.reportNowAsAdditional {
		return false
	}

	if f.mIncludeOnlyMulticastBefore.IsZero() && f.mIncludeOnlyMulticastBefore.Before(record.lastMulticastTime) {
		return false
	}

	if f.mReplyFilter != nil && !f.mReplyFilter.Accept(record.Responder.GetQType(), record.Responder.GetQClass(), record.Responder.GetQName()) {
		return false
	}

	return true
}

func (f *QueryResponderRecordFilter) SetIncludeAdditionalRepliesOnly(b bool) *QueryResponderRecordFilter {
	f.mIncludeAdditionalRepliesOnly = b
	return f
}
