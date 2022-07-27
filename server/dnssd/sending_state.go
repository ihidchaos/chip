package dnssd

import (
	"net"
	"net/netip"
)

const (
	ResourceType_Answer = iota
	ResourceType_Authority
	ResourceType_Additional
)

type ResponseSendingState struct {
	mQuery        *QueryData
	mMessageId    uint16
	mResourceType int
	mSendError    error
	mInterfaceId  net.Interface
	mSrcAddrPort  netip.AddrPort
	mDestAddrPort netip.AddrPort
}

func (s *ResponseSendingState) Reset(messageId uint16, query *QueryData, src, dest netip.AddrPort, interfaceId net.Interface) {
	s.mMessageId = messageId
	s.mQuery = query
	s.mInterfaceId = interfaceId
	s.mResourceType = ResourceType_Answer
	s.mSrcAddrPort = src
	s.mDestAddrPort = dest
}

func (s *ResponseSendingState) SetError(err error) {
	s.mSendError = err
}

func (s *ResponseSendingState) SendUnicast() bool {
	if s.mQuery == nil {
		return false
	}
	return s.mQuery.RequestedUnicastAnswer() || s.mDestAddrPort.Port() != MdnsPort
}

func (s *ResponseSendingState) GetMessageId() uint16 {
	return s.mMessageId
}

func (s *ResponseSendingState) GetError() error {
	return s.mSendError
}

func (s *ResponseSendingState) GetSourceAddress() netip.Addr {
	return s.mSrcAddrPort.Addr()
}

func (s *ResponseSendingState) SetResourceType(additional int) {
	s.mResourceType = additional
}

func (s *ResponseSendingState) SetSourceAddrPort(addr string, port uint16) {
	a, err := netip.ParseAddr(addr)
	if err != nil {
		return
	}
	s.mSrcAddrPort = netip.AddrPortFrom(a, port)
}

func (s *ResponseSendingState) IncludeQuery() bool {
	return s.mSrcAddrPort.Port() != MdnsPort
}

func (s *ResponseSendingState) GetQuery() *QueryData {
	return s.mQuery
}

func (s *ResponseSendingState) GetResourceType() int {
	return s.mResourceType
}
