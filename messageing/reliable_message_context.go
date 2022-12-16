package messageing

import (
	"github.com/galenliu/chip/lib/bitflags"
	"time"
)

const (
	fInitiator                uint16 = 1 << iota
	fResponseExpected         uint16 = 1 << 1
	fAutoRequestAck           uint16 = 1 << 2
	fMessageNotAcked          uint16 = 1 << 3
	fAckPending               uint16 = 1 << 4
	fAckMessageCounterIsValid uint16 = 1 << 5
	fWillSendMessage          uint16 = 1 << 6
	fClosed                   uint16 = 1 << 7
	fActiveMode               uint16 = 1 << 8
	fEphemeralExchange        uint16 = 1 << 9
	fIgnoreSessionRelease     uint16 = 1 << 10
)

type ReliableMessageContext struct {
	mFlags                        bitflags.Flags[uint16]
	mNextAckTime                  time.Time
	mPendingPeerAckMessageCounter uint32
	mExchangeContext              *ExchangeContext
}

func NewReliableMessageContext(ec *ExchangeContext) *ReliableMessageContext {
	return &ReliableMessageContext{
		mNextAckTime:                  time.Time{},
		mPendingPeerAckMessageCounter: 0,
		mExchangeContext:              ec,
	}
}

func (c *ReliableMessageContext) AutoRequestAck() bool {
	return c.mFlags.Has(fAutoRequestAck)

}

func (c *ReliableMessageContext) isAckPending() bool {
	return c.mFlags.Has(fAckPending)
}

func (c *ReliableMessageContext) isMessageNotAcked() bool {
	return c.mFlags.Has(fEphemeralExchange)
}

func (c *ReliableMessageContext) HasPiggybackAckPending() bool {
	return c.mFlags.Has(fAckMessageCounterIsValid)
}

func (c *ReliableMessageContext) IsRequestingActiveMode() bool {

	return c.mFlags.Has(fActiveMode)
}

func (c *ReliableMessageContext) SetAutoRequestAck(autoReqAck bool) {
	c.mFlags.Set(autoReqAck, fAutoRequestAck)
}

func (c *ReliableMessageContext) SetAckPending(inAckPending bool) {
	c.mFlags.Set(inAckPending, fAckPending)
}

func (c *ReliableMessageContext) SetMessageNotAcked(messageNotAcked bool) {
	c.mFlags.Set(messageNotAcked, fMessageNotAcked)
}

func (c *ReliableMessageContext) SetRequestingActiveMode(activeMode bool) {

	c.mFlags.Set(activeMode, fActiveMode)
}

func (c *ReliableMessageContext) IsEphemeralExchange() bool {
	return c.mFlags.Has(fEphemeralExchange)
}

func (c *ReliableMessageContext) HandleRcvdAck(ackMessageCounter uint32) {
	if !c.ReliableMessageMgr().CheckAndRemRetransTable(c, ackMessageCounter) {

	}
}

func (c *ReliableMessageContext) HandleNeedsAck(messageCounter, flags uint32) {
	//TODO implement me
	panic("implement me")
}

func (c *ReliableMessageContext) FlushAcks() error {
	//TODO implement me
	panic("implement me")
}

func (c *ReliableMessageContext) ShouldIgnoreSessionRelease() bool {
	return c.mFlags.Has(fIgnoreSessionRelease)
}

func (c *ReliableMessageContext) ReliableMessageMgr() *ReliableMessageMgr {
	return c.mExchangeContext.mExchangeMgr.mReliableMessageMgr
}

func (c *ReliableMessageContext) TakePendingPeerAckMessageCounter() uint32 {
	c.SetAckPending(false)
	return c.mPendingPeerAckMessageCounter
}
