package messageing

import (
	"github.com/galenliu/chip/lib/bitflags"
	"time"
)

const (
	kFlagInitiator                uint16 = 1 << iota
	kFlagResponseExpected         uint16 = 1 << 1
	kFlagAutoRequestAck           uint16 = 1 << 2
	kFlagMessageNotAcked          uint16 = 1 << 3
	kFlagAckPending               uint16 = 1 << 4
	kFlagAckMessageCounterIsValid uint16 = 1 << 5
	kFlagWillSendMessage          uint16 = 1 << 6
	kFlagClosed                   uint16 = 1 << 7
	kFlagActiveMode               uint16 = 1 << 8
	kFlagEphemeralExchange        uint16 = 1 << 9
	kFlagIgnoreSessionRelease     uint16 = 1 << 10
)

type ReliableMessageContext struct {
	mFlags                        bitflags.Flags[uint16]
	mNextAckTime                  time.Time
	mPendingPeerAckMessageCounter uint32
}

func NewReliableMessageContext() ReliableMessageContext {
	return ReliableMessageContext{
		mNextAckTime:                  time.Time{},
		mPendingPeerAckMessageCounter: 0,
	}
}

func (c *ReliableMessageContext) AutoRequestAck() bool {
	return c.mFlags.Has(kFlagAutoRequestAck)

}

func (c *ReliableMessageContext) IsAckPending() bool {

	return c.mFlags.Has(kFlagAckPending)
}

func (c *ReliableMessageContext) IsMessageNotAcked() bool {
	return c.mFlags.Has(kFlagEphemeralExchange)
}

func (c *ReliableMessageContext) HasPiggybackAckPending() bool {
	return c.mFlags.Has(kFlagAckMessageCounterIsValid)
}

func (c *ReliableMessageContext) IsRequestingActiveMode() bool {

	return c.mFlags.Has(kFlagActiveMode)
}

func (c *ReliableMessageContext) SetAutoRequestAck(autoReqAck bool) {
	c.mFlags.Sets(autoReqAck, kFlagAutoRequestAck)
}

func (c *ReliableMessageContext) SetAckPending(inAckPending bool) {
	c.mFlags.Sets(inAckPending, kFlagAckPending)
}

func (c *ReliableMessageContext) SetMessageNotAcked(messageNotAcked bool) {
	c.mFlags.Sets(messageNotAcked, kFlagMessageNotAcked)
}

func (c *ReliableMessageContext) SetRequestingActiveMode(activeMode bool) {

	c.mFlags.Sets(activeMode, kFlagActiveMode)
}

func (c *ReliableMessageContext) IsEphemeralExchange() bool {
	return c.mFlags.Has(kFlagEphemeralExchange)
}

func (c *ReliableMessageContext) HandleRcvdAck(ackCounter uint32) {
	//TODO implement me
	panic("implement me")
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
	return c.mFlags.Has(kFlagIgnoreSessionRelease)
}
