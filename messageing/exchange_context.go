package messageing

import (
	"fmt"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/platform/system"
	"github.com/galenliu/chip/protocols"
	"github.com/galenliu/chip/protocols/secure_channel"
	log "golang.org/x/exp/slog"
	"time"
)

type ExchangeSessionHolder struct {
	*transport.SessionHolderWithDelegate
}

// GrabExpiredSession 抓取过期会话
func (s *ExchangeSessionHolder) GrabExpiredSession(session *transport.SessionHandle) {
	secureSession, ok := session.Session.(*transport.SecureSession)
	if ok {
		if secureSession.IsPendingEviction() {
			s.GrabUnchecked(session)
		}
	}
}

func NewExchangeSessionHolder(delegate *ExchangeContext) *ExchangeSessionHolder {
	return &ExchangeSessionHolder{
		SessionHolderWithDelegate: transport.NewSessionHolderWithDelegateImpl(delegate),
	}
}

type ExchangeContext struct {
	ReliableMessageContext
	mExchangeId      uint16
	mExchangeMgr     *ExchangeManager
	mDispatch        ExchangeMessageDispatchBase
	mSession         *ExchangeSessionHolder
	mDelegate        ExchangeDelegate
	mFlags           uint16
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
	var flags uint16 = 0
	flags = lib.SetFlag(initiator, flags, kFlagInitiator)
	flags = lib.SetFlag(isEphemeralExchange, flags, kFlagEphemeralExchange)
	ec := &ExchangeContext{
		ReliableMessageContext: NewReliableMessageContext(),
		mExchangeId:            exchangeId,
		mExchangeMgr:           em,
		mDelegate:              delegate,
		mFlags:                 flags,
	}
	ec.mDispatch = ec.GetMessageDispatch(isEphemeralExchange, delegate)
	ec.ReferenceCounted = lib.NewReferenceCounted(1, ec)
	ec.mSession = NewExchangeSessionHolder(ec)
	ec.mSession.Grad(session)
	ec.SetResponseTimeout(time.Duration(0))

	if ec.IsInitiator() && !isEphemeralExchange {
		ec.WillSendMessage()
	}
	ec.SetAckPending(false)

	return ec
}

// IsInitiator 是否是通信的发起方
func (c *ExchangeContext) IsInitiator() bool {
	return c.mFlags&kFlagInitiator != 0
}

func (c *ExchangeContext) IsEncryptionRequired() bool {
	return c.mDispatch.IsEncryptionRequired()
}

func (c *ExchangeContext) IsGroupExchangeContext() bool {
	return c.mSession.IsGroupSession()
}

func (c *ExchangeContext) SendMessage(protocolId *protocols.Id, msgType uint8, r2 []byte, response uint16) error {

	//isStandaloneAck := protocolId == protocols.StandardSecureChannelProtocolId && msgType == StandaloneAck

	return nil
}

func (c *ExchangeContext) MatchExchange(session *transport.SessionHandle, packetHeader *raw.PacketHeader, payloadHeader *raw.PayloadHeader) bool {
	return (c.mExchangeId == payloadHeader.ExchangeId()) &&
		(c.mSession.Contains(session)) &&
		(c.IsEncryptionRequired() == packetHeader.IsEncrypted()) &&
		(payloadHeader.IsInitiator() != c.IsInitiator())
}

func (c *ExchangeContext) HandleMessage(messageCounter uint32, payloadHeader *raw.PayloadHeader, flags uint32, buf *system.PacketBufferHandle) error {

	c.Retain()
	isStandaloneAck := payloadHeader.HasMessageType(uint8(secure_channel.StandaloneAck))
	isDuplicate := lib.HasFlags(flags, fDuplicateMessage)

	defer func() {
		if (isDuplicate || isStandaloneAck) && c.mDelegate != nil {
			return
		}
		c.messageHandled()
	}()

	if c.mDispatch.IsReliableTransmissionAllowed() && !c.IsGroupExchangeContext() {
		if lib.HasFlags(flags, fDuplicateMessage) &&
			payloadHeader.IsAckMsg() &&
			payloadHeader.AckMessageCounter().Present() {
			c.HandleRcvdAck(payloadHeader.AckMessageCounter().MustGet())
		}
		if payloadHeader.NeedsAck() {
			c.HandleNeedsAck(messageCounter, flags)
		}
	}
	if c.isAckPending() && c.mDelegate != nil {
		return c.FlushAcks()
	}

	if isStandaloneAck || isDuplicate || c.IsEphemeralExchange() {
		return nil
	}

	if c.IsMessageNotAcked() {
		log.Info("ExchangeManager Dropping message without piggyback ack when we are waiting for an ack.")
		return lib.MATTER_ERROR_INCORRECT_STATE
	}
	c.CancelResponseTimer()
	c.SetResponseExpected(false)
	if c.mDelegate != nil && c.mDispatch.MessagePermitted(payloadHeader.ProtocolID(), payloadHeader.MessageType()) {
		return c.mDelegate.OnMessageReceived(c, payloadHeader, buf)
	}
	DefaultOnMessageReceived(c, payloadHeader.ProtocolID(), payloadHeader.MessageType(), messageCounter, buf)
	return nil
}

func (c *ExchangeContext) SetDelegate(delegate ExchangeDelegate) {
	c.mDelegate = delegate
}

func (c *ExchangeContext) GetMessageDispatch(isEphemeralExchange bool, delegate ExchangeDelegate) ExchangeMessageDispatchBase {
	if isEphemeralExchange {
		return EphemeralExchangeDispatchInstance()
	}
	if c.mDelegate != nil {
		return c.mDelegate.GetMessageDispatch()
	}
	return ApplicationExchangeDispatchInstance()
}

func (c *ExchangeContext) Delegate() ExchangeDelegate {
	return c.mDelegate
}

func (c *ExchangeContext) SetResponseTimeout(timeout time.Duration) {
	c.mResponseTimeout = timeout
}

func (c *ExchangeContext) Released() {
	c.mExchangeMgr.ReleaseContext(c)
}

func (c *ExchangeContext) DoClose(clearRetransTable bool) {
	if lib.HasFlags(c.mFlags, kFlagClosed) {
		return
	}
	c.mFlags = lib.SetFlags(c.mFlags, kFlagClosed)

	if c.mDelegate != nil {
		c.mDelegate.OnExchangeClosing(c)
	}
	c.mDelegate = nil
	_ = c.FlushAcks()

	if clearRetransTable {
		c.mExchangeMgr.ReliableMessageMgr().ClearRetransTable(c)
	}
	c.CancelResponseTimer()
}

func (c *ExchangeContext) Close() {
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

func (c *ExchangeContext) sendMessage(id protocols.Id) error {

	//c.mDispatch.SendMessage(c.ExchangeMgr().SessionManagerBase(),c.mSession.)
	return nil
}

func (c *ExchangeContext) HasSessionHandle() bool {
	return c.mSession != nil
}

func (c *ExchangeContext) GetSessionHandle() *transport.SessionHandle {
	return c.mSession.Get()
}

func (c *ExchangeContext) CancelResponseTimer() {
	systemLayer := c.mExchangeMgr.SessionManager().SystemLayer()
	if systemLayer == nil {
		return
	}
	systemLayer.CancelTimer(c.HandleResponseTimeout, c)
}

func (c *ExchangeContext) HandleResponseTimeout(layer system.Layer, aAppState any) {
	ec, ok := aAppState.(*ExchangeContext)
	if !ok {
		return
	}
	ec.NotifyResponseTimeout(true)
}

func (c *ExchangeContext) NotifyResponseTimeout(b bool) {

}

func (c *ExchangeContext) LogValue() log.Value {
	return log.GroupValue(
		log.Int("id", int(c.mExchangeId)),
		log.String("delegate", fmt.Sprintf("%p", c.Delegate())),
	)
}

func (c *ExchangeContext) messageHandled() {
	if lib.HasFlags(c.mFlags, kFlagClosed) || c.isResponseExpected() || c.isSendExpected() {
		return
	}
	c.Close()
}

func (c *ExchangeContext) isResponseExpected() bool {
	return lib.HasFlags(c.mFlags, kFlagResponseExpected)
}

func (c *ExchangeContext) isSendExpected() bool {
	return lib.HasFlags(c.mFlags, kFlagWillSendMessage)
}

func (c *ExchangeContext) isAckPending() bool {
	return lib.HasFlags(c.mFlags, kFlagAckPending)
}

func (c *ExchangeContext) SetResponseExpected(b bool) {
	c.mFlags = lib.SetFlag(b, c.mFlags, kFlagResponseExpected)
}

func (c *ExchangeContext) WillSendMessage() {
	c.mFlags = lib.SetFlag(true, c.mFlags, kFlagWillSendMessage)
}

func DefaultOnMessageReceived(c *ExchangeContext, id *protocols.Id, messageType uint8, messageCounter uint32, payload *system.PacketBufferHandle) {
	log.Error("ExchangeManager Dropping unexpected message of type", lib.MATTER_ERROR_INVALID_MESSAGE_TYPE,
		"MessageType", messageType,
		"protocolId", id,
		"MessageCounter", messageCounter,
		"exchange", c)
}
