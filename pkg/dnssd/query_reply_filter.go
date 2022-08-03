package dnssd

import (
	"github.com/galenliu/chip/pkg/dnssd/responder"
	"github.com/miekg/dns"
)

type QueryReplyFilter struct {
	mIgnoreNameMatch        bool
	mSendingAdditionalItems bool
	mQueryData              *QueryData
	responder.ReplyFilter
}

func NewQueryReplyFilter(q *QueryData) *QueryReplyFilter {
	return &QueryReplyFilter{
		mIgnoreNameMatch:        false,
		mSendingAdditionalItems: false,
		mQueryData:              q,
	}
}

func (f *QueryReplyFilter) Accept(qType, qClass uint16, fName string) bool {
	if !f.acceptableQueryType(qType) {
		return false
	}

	if !f.acceptableQueryClass(qClass) {
		return false
	}
	return f.acceptablePath(fName)
}

func (f *QueryReplyFilter) acceptableQueryType(qType uint16) bool {
	if f.mSendingAdditionalItems {
		return true
	}
	return (f.mQueryData.GetType() == dns.TypeANY) || (f.mQueryData.GetType() == qType)
}

func (f *QueryReplyFilter) acceptableQueryClass(qClass uint16) bool {
	return (f.mQueryData.GetClass() == dns.ClassANY) || (f.mQueryData.GetClass() == qClass)
}

func (f *QueryReplyFilter) acceptablePath(qName string) bool {
	if f.mIgnoreNameMatch || f.mQueryData.IsInternalBroadcast() {
		return true
	}
	return f.mQueryData.GetName() == qName
}

func (f *QueryReplyFilter) SetIgnoreNameMatch(b bool) *QueryReplyFilter {
	f.mIgnoreNameMatch = b
	return f
}

func (f *QueryReplyFilter) SetSendingAdditionalItems(b bool) *QueryReplyFilter {
	f.mSendingAdditionalItems = b
	return f
}
