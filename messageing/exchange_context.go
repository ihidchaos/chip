package messageing

import (
	"fmt"
	"github.com/galenliu/chip"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/lib/bitflags"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/messageing/transport/session"
	"github.com/galenliu/chip/platform/system"
	"github.com/galenliu/chip/protocols"
	log "golang.org/x/exp/slog"
	"time"
)

var cancelChan chan any

type ExchangeHandle struct {
	*ExchangeContext
}

type ExchangeSessionHolder struct {
	*transport.SessionHolderWithDelegate
}

// GrabExpiredSession 抓取过期会话
func (s *ExchangeSessionHolder) GrabExpiredSession(ss *transport.SessionHandle) {
	secureSession, ok := ss.Session.(*session.Secure)
	if ok {
		if secureSession.IsPendingEviction() {
			s.GrabUnchecked(ss)
		}
	}
}

func NewExchangeSessionHolder(delegate *ExchangeContext) *ExchangeSessionHolder {
	return &ExchangeSessionHolder{
		SessionHolderWithDelegate: transport.NewSessionHolderWithDelegateImpl(delegate),
	}
}

type ExchangeContext struct {
	*ReliableMessageContext
	mExchangeId      uint16
	mExchangeMgr     *ExchangeManager
	mDispatch        ExchangeMessageDispatch
	mSession         *ExchangeSessionHolder
	mDelegate        ExchangeDelegate
	mFlags           bitflags.Flags[uint16]
	mResponseTimeout time.Duration
	*lib.ReferenceCounted
}

func NewExchangeContext(
	em *ExchangeManager,
	exchangeId uint16,
	session *transport.SessionHandle,
	initiator bool,
	delegate ExchangeDelegate,
	isEphemeralExchange bool,
) *ExchangeContext {
	ec := &ExchangeContext{
		mExchangeId:  exchangeId,
		mExchangeMgr: em,
		mDelegate:    delegate,
		mFlags:       bitflags.Flags[uint16]{},
	}
	ec.ReliableMessageContext = NewReliableMessageContext(ec)
	ec.mFlags.Set(initiator, fInitiator)
	ec.mFlags.Set(isEphemeralExchange, fEphemeralExchange)
	ec.mDispatch = ec.GetMessageDispatch(isEphemeralExchange, delegate)
	ec.ReferenceCounted = lib.NewReferenceCounted(1, ec)
	ec.mSession = NewExchangeSessionHolder(ec)
	ec.mSession.Grad(session)
	ec.SetResponseTimeout(time.Duration(0))
	if ec.isInitiator() && !isEphemeralExchange {
		ec.WillSendMessage()
	}
	ec.SetAckPending(false)
	ec.SetAutoRequestAck(!session.IsGroup())
	return ec
}

// isInitiator 是否是通信的发起方
func (c *ExchangeContext) isInitiator() bool {
	return c.mFlags.Has(fInitiator)
}

func (c *ExchangeContext) isResponseExpected() bool {
	return c.mFlags.Has(fResponseExpected)
}

func (c *ExchangeContext) setResponseExpected(inResponseExpected bool) {
	c.mFlags.Set(inResponseExpected, fResponseExpected)

}

func (c *ExchangeContext) UseSuggestedResponseTimeout(applicationProcessingTimeout time.Duration) {
	c.SetResponseTimeout(c.mSession.ComputeRoundTripTimeout(applicationProcessingTimeout))
}

func (c *ExchangeContext) SetResponseTimeout(timeout time.Duration) {
	c.mResponseTimeout = timeout
}

func (c *ExchangeContext) isEncryptionRequired() bool {
	return c.mDispatch.IsEncryptionRequired()
}

func (c *ExchangeContext) isGroupExchangeContext() bool {
	return c.mSession != nil && c.mSession.IsGroup()
}

func (c *ExchangeContext) SendMessage(msgType MessageType, msg []byte, sendFlags bitflags.Flags[uint16]) (err error) {
	isStandaloneAck := lib.IsStandaloneAck(msgType.MessageType())

	if c.mExchangeMgr == nil {
		return chip.New(chip.ErrorInternal, "ExchangeContext", "ExchangeMgr nil")
	}
	if c.mSession == nil {
		return chip.New(chip.ErrorConnectionAborted)
	}
	reliableTransmissionRequested := c.SessionHandle().RequireMRP() && sendFlags.Has(fNoAutoRequestAck) && !c.isGroupExchangeContext()

	if sendFlags.Has(fExpectResponse) && !c.isGroupExchangeContext() {
		if c.isResponseExpected() {
			return chip.ErrorIncorrectState
		}
		c.setResponseExpected(true)
		if c.mResponseTimeout > 0 {
			if err = c.startResponseTimer(time.After(c.mResponseTimeout)); err != nil {
				c.setResponseExpected(false)
				return err
			}
		}
	}
	if c.isGroupExchangeContext() && !c.isInitiator() {
		return chip.ErrorInternal
	}

	err = c.mDispatch.SendMessage(c.mExchangeMgr.mSessionManager,
		c.mSession.SessionHandler(),
		c.mExchangeId, c.isInitiator(),
		c.ReliableMessageContext,
		reliableTransmissionRequested,
		msgType.ProtocolId(),
		msgType.MessageType(),
		msg)

	if err != nil && c.isResponseExpected() {
		c.cancelResponseTimer()
		c.setResponseExpected(false)
	}
	if err == nil && !isStandaloneAck {
		c.mFlags.Clear(fWillSendMessage)
		c.messageHandled()
	}
	return err
}

func (c *ExchangeContext) MatchExchange(session *transport.SessionHandle,
	packetHeader *raw.PacketHeader,
	payloadHeader *raw.PayloadHeader) bool {
	return (c.mExchangeId == payloadHeader.ExchangeId) &&
		(c.mSession.Contains(session)) &&
		(c.isEncryptionRequired() == packetHeader.IsEncrypted()) &&
		(payloadHeader.IsInitiator() != c.isInitiator())
}

func (c *ExchangeContext) HandleMessage(messageCounter uint32, payloadHeader *raw.PayloadHeader, f uint32, buf *system.PacketBufferHandle) error {
	msgFlags := bitflags.Some(f)
	isDuplicate := msgFlags.Has(fDuplicateMessage)
	var standaloneAck uint8 = 0x10
	isStandaloneAck := payloadHeader.HasMessageType(standaloneAck)
	defer func() {
		if (isDuplicate || isStandaloneAck) && c.mDelegate != nil {
			return
		}
		c.messageHandled()
	}()

	if c.mDispatch.IsReliableTransmissionAllowed() && !c.isGroupExchangeContext() {
		if !msgFlags.Has(fDuplicateMessage) &&
			payloadHeader.IsAckMsg() &&
			payloadHeader.AckMessageCounter.IsSome() {
			c.HandleRcvdAck(payloadHeader.AckMessageCounter.Unwrap())
		}
		if payloadHeader.NeedsAck() {
			c.HandleNeedsAck(messageCounter, f)
		}
	}
	if c.isAckPending() && c.mDelegate != nil {
		err := c.FlushAcks()
		if err != nil {
			return err
		}
	}

	if isStandaloneAck || isDuplicate || c.IsEphemeralExchange() {
		return nil
	}

	if c.isMessageNotAcked() {
		log.Info("ExchangeManager Dropping message without piggyback ack when we are waiting for an ack.")
		return chip.ErrorIncorrectState
	}

	c.cancelResponseTimer()

	c.setResponseExpected(false)
	if c.mDelegate != nil && c.mDispatch.MessagePermitted(payloadHeader.ProtocolId, payloadHeader.MessageType) {
		return c.mDelegate.OnMessageReceived(c, payloadHeader, buf)
	}
	DefaultOnMessageReceived(c, payloadHeader.ProtocolId, payloadHeader.MessageType, messageCounter, buf)
	return nil
}

func (c *ExchangeContext) SetDelegate(delegate ExchangeDelegate) {
	c.mDelegate = delegate
}

func (c *ExchangeContext) ExchangeId() uint16 {
	return c.mExchangeId
}

func (c *ExchangeContext) GetMessageDispatch(isEphemeralExchange bool,
	delegate ExchangeDelegate) ExchangeMessageDispatch {
	if isEphemeralExchange {
		return DefaultEphemeralDispatch()
	}
	if delegate != nil {
		return delegate.GetMessageDispatch()
	}
	return DefaultApplicationExchangeDispatch()
}

func (c *ExchangeContext) Delegate() ExchangeDelegate {
	return c.mDelegate
}

func (c *ExchangeContext) Released() {
	c.mExchangeMgr.ReleaseContext(c)
}

func (c *ExchangeContext) DoClose(clearRetransTable bool) {

	if c.mFlags.Has(fClosed) {
		return
	}
	c.mFlags.Set(true, fClosed)

	if c.mDelegate != nil {
		c.mDelegate.OnExchangeClosing(c)
	}
	c.mDelegate = nil
	_ = c.FlushAcks()

	if clearRetransTable {
		c.mExchangeMgr.ReliableMessageMgr().ClearRetransTable(c.ReliableMessageContext)
	}
	c.cancelResponseTimer()
}

func (c *ExchangeContext) close() {
	c.DoClose(false)
	c.Release()
}

func (c *ExchangeContext) ExchangeMgr() *ExchangeManager {
	return c.mExchangeMgr
}

func (c *ExchangeContext) OnSessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (c *ExchangeContext) HasSessionHandle() bool {
	return c.mSession != nil
}

func (c *ExchangeContext) SessionHandle() *transport.SessionHandle {
	return c.mSession.SessionHandler()
}

func (c *ExchangeContext) cancelResponseTimer() {
	if cancelChan != nil {
		select {
		case cancelChan <- true:
		}
	}
}

func (c *ExchangeContext) startResponseTimer(timOutChan <-chan time.Time) error {
	if cancelChan == nil {
		cancelChan = make(chan any)
	}
	go func() {
		select {
		case <-timOutChan:
			c.handleResponseTimeout()
			return
		case <-cancelChan:
			return
		}
	}()
	return nil
}

func (c *ExchangeContext) handleResponseTimeout() {
	c.notifyResponseTimeout(true)
}

func (c *ExchangeContext) notifyResponseTimeout(aCloseIfNeeded bool) {
	c.setResponseExpected(false)
	if c.mSession != nil {
		if c.mSession.IsSecure() && c.mSession.Session.(*session.Secure).IsCASESession() {
			c.mSession.Session.(*session.Secure).MarkAsDefunct()
		}
	}

}

func (c *ExchangeContext) LogValue() log.Value {
	return log.GroupValue(
		log.Int("id", int(c.mExchangeId)),
		log.String("delegate", fmt.Sprintf("%p", c.Delegate())),
	)
}

func (c *ExchangeContext) messageHandled() {
	if c.mFlags.Has(fClosed) || c.isResponseExpected() || c.isSendExpected() {
		return
	}
	c.close()
}

func (c *ExchangeContext) isSendExpected() bool {
	return c.mFlags.Has(fWillSendMessage)
}

func (c *ExchangeContext) isAckPending() bool {
	return c.mFlags.Has(fAckPending)

}

func (c *ExchangeContext) WillSendMessage() {
	c.mFlags.Set(true, fWillSendMessage)
}

func DefaultOnMessageReceived(c *ExchangeContext, id protocols.Id, messageType uint8, messageCounter uint32, payload *system.PacketBufferHandle) {
	log.Error("ExchangeManager Dropping unexpected message of type", chip.ErrorInvalidMessageType,
		"MessageType", messageType,
		"protocolId", id,
		"MessageCounterBase", messageCounter,
		"exchange", c)
}
