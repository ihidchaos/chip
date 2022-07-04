package parameters

import (
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/core"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/pkg"
	"net"
)

type Mac struct {
	mac string
}

type BaseAdvertisingParams struct {
	mPort           uint16
	mMac            string
	mEnableIPv4     bool
	mInterfaceId    net.Interface
	mMRPConfig      *messageing.ReliableMessageProtocolConfig
	mTcpSupported   *bool
	mLocalMRPConfig *messageing.ReliableMessageProtocolConfig
}

type CommissionAdvertisingParameters struct {
	*BaseAdvertisingParams
	mVendorId          *uint16 //供应商口称
	mProductId         *uint16 //产品ID
	mDeviceType        *int32  //设备类型
	mPairingHint       *uint16 //设备配提示
	mPairingInstr      string  //设备配对指南
	mDeviceName        string  //设备名称
	mMode              config.CommssionAdvertiseMode
	mCommissioningMode config.CommissioningMode
	mPeerId            *core.PeerId
	mLongDiscriminator uint16
	mRotatingId        string
}

type OperationalAdvertisingParameters struct {
	*BaseAdvertisingParams
	mPeerId core.PeerId
}

func (o *OperationalAdvertisingParameters) SetPeerId(peerId core.PeerId) {
	o.mPeerId = peerId
}

func (o *OperationalAdvertisingParameters) GetCompressedFabricId() core.CompressedFabricId {
	return o.mPeerId.GetCompressedFabricId()
}

func (o *OperationalAdvertisingParameters) GetPeerId() core.PeerId {
	return o.mPeerId
}

func (c *CommissionAdvertisingParameters) SetCommissioningMode(mode config.CommissioningMode) {
	c.mCommissioningMode = mode
}

func (c *CommissionAdvertisingParameters) GetCommissioningMode() config.CommissioningMode {
	return c.mCommissioningMode
}

func (c *CommissionAdvertisingParameters) SetCommissionAdvertiseMode(mode config.CommssionAdvertiseMode) {
	c.mMode = mode
}

func (c *CommissionAdvertisingParameters) GetCommissionAdvertiseMode() config.CommssionAdvertiseMode {
	return c.mMode
}

func (c *CommissionAdvertisingParameters) SetVendorId(id uint16) {
	c.mVendorId = &id
}

func (c *CommissionAdvertisingParameters) SetProductId(id uint16) {
	c.mProductId = &id
}

func (c *CommissionAdvertisingParameters) SetDeviceType(t int32) {
	c.mDeviceType = &t
}

func (c *CommissionAdvertisingParameters) SetDeviceName(name string) {
	c.mDeviceName = name
}

func (c *CommissionAdvertisingParameters) SetTcpSupported(b bool) {
	c.mTcpSupported = &b
}

func (c *CommissionAdvertisingParameters) SetPairingHint(value uint16) {
	c.mPairingHint = &value
}

func (c *CommissionAdvertisingParameters) SetPairingInstruction(ist string) {
	c.mPairingInstr = ist
}

func (c *CommissionAdvertisingParameters) SetMRPConfig(config *messageing.ReliableMessageProtocolConfig) {
	c.mMRPConfig = config
}

func (c *CommissionAdvertisingParameters) GetVendorId() *uint16 {
	return c.mVendorId
}

func (c *CommissionAdvertisingParameters) GetDeviceType() *int32 {
	return c.mDeviceType
}

func (c *CommissionAdvertisingParameters) GetProductId() *uint16 {
	return c.mProductId
}

func (c *CommissionAdvertisingParameters) GetDeviceName() string {
	return c.mDeviceName
}

func (c *CommissionAdvertisingParameters) GetLongDiscriminator() uint16 {
	return c.mLongDiscriminator
}

func (c *CommissionAdvertisingParameters) SetLongDiscriminator(discriminator uint16) *CommissionAdvertisingParameters {
	c.mLongDiscriminator = discriminator
	return c
}

func (c *CommissionAdvertisingParameters) GetRotatingDeviceId() string {
	return c.mRotatingId
}

func (c *CommissionAdvertisingParameters) GetPairingHint() *uint16 {
	return c.mPairingHint
}

func (c *CommissionAdvertisingParameters) GetPairingInstruction() string {
	return c.mPairingInstr
}

func (b *BaseAdvertisingParams) IsIPv4Enabled() bool {
	return b.mEnableIPv4
}

func (b *BaseAdvertisingParams) SetPort(port uint16) {
	b.mPort = port
}

func (b *BaseAdvertisingParams) GetPort() uint16 {
	return b.mPort
}

func (b *BaseAdvertisingParams) SetMaC(mac string) {
	b.mMac = mac
}

func (b *BaseAdvertisingParams) GetMac() string {
	if b.mMac == "" {
		b.mMac = pkg.Mac48Address(pkg.RandHex())
	}
	return b.mMac
}

func (b *BaseAdvertisingParams) GetUUID() string {
	if b.mMac == "" {
		b.mMac = pkg.Mac48Address(pkg.RandHex())
	}
	return b.mMac
}

func (b *BaseAdvertisingParams) EnableIpV4(enable bool) {
	b.mEnableIPv4 = enable
}

func (b *BaseAdvertisingParams) GetLocalMRPConfig() interface{} {
	return b.mLocalMRPConfig
}

func (b *BaseAdvertisingParams) GetTcpSupported() *bool {
	return b.mTcpSupported
}
