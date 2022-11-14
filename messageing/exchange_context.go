package messageing

import (
	"encoding/json"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/platform/system"
	"github.com/galenliu/chip/protocols"
	log "golang.org/x/exp/slog"
)

type ExchangeSessionHolder struct {
	*transport.SessionHolderWithDelegate
}

func NewExchangeSessionHolder(delegate *ExchangeContext) *ExchangeSessionHolder {
	return &ExchangeSessionHolder{
		SessionHolderWithDelegate: transport.NewSessionHolderWithDelegateImpl(delegate),
	}
}

type ExchangeContext struct {
	ReliableMessageContext
	mExchangeId  uint16
	mExchangeMgr *ExchangeManager
	mDispatch    ExchangeMessageDispatch
	mSession     *ExchangeSessionHolder
	mDelegate    ExchangeDelegate
	mFlags       uint16

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
		ReliableMessageContext: nil,
		mExchangeId:            exchangeId,
		mExchangeMgr:           em,
		mDispatch:              GetMessageDispatch(isEphemeralExchange, delegate),
		mDelegate:              delegate,
		mFlags:                 flags,
	}
	ec.mSession = NewExchangeSessionHolder(ec)
	ec.ReferenceCounted = lib.NewReferenceCounted(1, ec)
	//ec.mSession.SessionHolderWithDelegate.Grad(session)
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

func (c *ExchangeContext) SendMessage(protocolId *protocols.Id, msgType MsgType, r2 []byte, response uint16) error {

	//isStandaloneAck := protocolId == protocols.StandardSecureChannelProtocolId && msgType == StandaloneAck

	return nil
}

func (c *ExchangeContext) MatchExchange(session *transport.SessionHandle, packetHeader *raw.PacketHeader, payloadHeader *raw.PayloadHeader) bool {
	return (c.mExchangeId == payloadHeader.ExchangeId()) &&
		(c.mSession.Contains(session)) &&
		(c.IsEncryptionRequired() == packetHeader.IsEncrypted()) &&
		(payloadHeader.IsInitiator() != c.IsInitiator())
}

func (c *ExchangeContext) HandleMessage(counter uint32, payloadHeader *raw.PayloadHeader, flags uint32, buf *system.PacketBufferHandle) error {

	//isStandaloneAck := payloadHeader.HasMessageType(uint8(StandaloneAck))
	//isDuplicate := lib.HasFlags(flags, transport.DuplicateMessageFlag)
	if c.mDelegate != nil {
		err := c.mDelegate.OnMessageReceived(c, payloadHeader, buf)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ExchangeContext) SetDelegate(delegate ExchangeDelegate) {

}

func GetMessageDispatch(isEphemeralExchange bool, delegate ExchangeDelegate) ExchangeMessageDispatch {
	if isEphemeralExchange {
		return EphemeralExchangeDispatchImpl{delegate: delegate}
	}
	return ExchangeMessageDispatchImpl{delegate: delegate}
}

func (c *ExchangeContext) GetDelegate() ExchangeDelegate {
	return c.mDelegate
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
	c.FlushAcks()

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

func (c *ExchangeContext) Marshall() string {
	data, _ := json.MarshalIndent(struct {
		ExchangeId uint16
		Flags      uint16
	}{
		ExchangeId: c.mExchangeId,
		Flags:      c.mFlags,
	}, "", "   ")
	return string(data)
}

func (c *ExchangeContext) FlushAcks() {

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
		log.Any("flags", c.mFlags),
		log.Any("exchangeId", c.mExchangeId),
	)
}
