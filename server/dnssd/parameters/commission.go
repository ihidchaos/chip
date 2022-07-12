package parameters

import (
	"fmt"
	"github.com/galenliu/chip/core"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/server/dnssd/costants/commssion_advertise_mode"
)

type Mac struct {
	mac string
}

type CommissionAdvertisingParameters struct {
	*BaseAdvertisingParams
	mVendorId           *uint16 //供应商口称
	mProductId          *uint16 //产品ID
	mDeviceType         *uint32 //设备类型
	mPairingHint        string  //设备配提示
	mPairingInstr       string  //设备配对指南
	mDeviceName         string  //设备名称
	mMode               CommssionAdvertiseMode.T
	mCommissioningMode  uint8
	mPeerId             *core.PeerId
	mShortDiscriminator uint8
	mLongDiscriminator  uint16
	mRotatingId         string
}

func (c CommissionAdvertisingParameters) Init() *CommissionAdvertisingParameters {
	c.BaseAdvertisingParams = BaseAdvertisingParams{}.Init()
	return &c
}

func (c *CommissionAdvertisingParameters) SetCommissioningMode(mode uint8) {
	c.mCommissioningMode = mode
}

func (c *CommissionAdvertisingParameters) GetCommissioningMode() uint8 {
	return c.mCommissioningMode
}

func (c *CommissionAdvertisingParameters) SetCommissionAdvertiseMode(mode CommssionAdvertiseMode.T) {
	c.mMode = mode
}

func (c *CommissionAdvertisingParameters) GetCommissionAdvertiseMode() CommssionAdvertiseMode.T {
	return c.mMode
}

func (c *CommissionAdvertisingParameters) SetVendorId(id uint16) {
	c.mVendorId = &id
}

func (c *CommissionAdvertisingParameters) SetProductId(id uint16) *CommissionAdvertisingParameters {
	c.mProductId = &id
	return c
}

func (c *CommissionAdvertisingParameters) SetDeviceType(t uint32) *CommissionAdvertisingParameters {
	c.mDeviceType = &t
	return c
}

func (c *CommissionAdvertisingParameters) SetDeviceName(name string) *CommissionAdvertisingParameters {
	c.mDeviceName = name
	return c
}

func (c *CommissionAdvertisingParameters) SetTcpSupported(b bool) *CommissionAdvertisingParameters {
	c.mTcpSupported = &b
	return c
}

func (c *CommissionAdvertisingParameters) SetPairingHint(value string) *CommissionAdvertisingParameters {
	c.mPairingHint = value
	return c
}

func (c *CommissionAdvertisingParameters) SetPairingInstruction(ist string) {
	c.mPairingInstr = ist
}

func (c *CommissionAdvertisingParameters) SetMRPConfig(config *messageing.ReliableMessageProtocolConfig) {
	c.mMRPConfig = config
}

func (c *CommissionAdvertisingParameters) GetVendorId() (uint16, error) {
	if c.mVendorId == nil {
		return 0, fmt.Errorf("vendor id not set")
	}
	return *c.mVendorId, nil
}

func (c *CommissionAdvertisingParameters) GetDeviceType() (t uint32, e error) {
	if c.mDeviceType == nil {
		e = fmt.Errorf("value not set")
		return
	}
	t = *c.mDeviceType
	return
}

func (c *CommissionAdvertisingParameters) GetProductId() (uint16, error) {
	if c.mProductId == nil {
		return 0, fmt.Errorf("product id not set")
	}
	return *c.mProductId, nil
}

func (c *CommissionAdvertisingParameters) GetDeviceName() string {
	return c.mDeviceName
}

func (c *CommissionAdvertisingParameters) GetLongDiscriminator() uint16 {
	return c.mLongDiscriminator
}

func (c *CommissionAdvertisingParameters) GetShortDiscriminator() uint8 {
	return c.mShortDiscriminator
}

func (c *CommissionAdvertisingParameters) SetLongDiscriminator(discriminator uint16) *CommissionAdvertisingParameters {
	c.mLongDiscriminator = discriminator
	return c
}

func (c *CommissionAdvertisingParameters) GetRotatingDeviceId() string {
	return c.mRotatingId
}

func (c *CommissionAdvertisingParameters) GetPairingHint() string {
	return c.mPairingHint
}

func (c *CommissionAdvertisingParameters) GetPairingInstruction() string {
	return c.mPairingInstr
}

func (c *CommissionAdvertisingParameters) SetShortDiscriminator(discriminator uint8) *CommissionAdvertisingParameters {
	c.mShortDiscriminator = discriminator
	return c
}

func (c *CommissionAdvertisingParameters) SetLocalMRPConfig(config *messageing.ReliableMessageProtocolConfig) *CommissionAdvertisingParameters {
	c.mLocalMRPConfig = config
	return c
}
