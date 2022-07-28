package params

import (
	"fmt"
	"github.com/galenliu/chip/device"
	"github.com/galenliu/chip/messageing"
)

type Mac struct {
	mac string
}

type CommissionAdvertisingParameters struct {
	BaseAdvertisingParams
	mVendorId           *uint16 //供应商口称
	mProductId          *uint16 //产品ID
	mDeviceType         *uint32 //设备类型
	mPairingHint        uint16  //设备配提示
	mPairingInstr       string  //设备配对指南
	mDeviceName         string  //设备名称
	mMode               uint16
	mCommissioningMode  int
	mPeerId             *device.PeerId
	mShortDiscriminator uint8
	mLongDiscriminator  uint16
	mRotatingId         string
}

func NewCommissionAdvertisingParameters() *CommissionAdvertisingParameters {
	return &CommissionAdvertisingParameters{
		BaseAdvertisingParams: BaseAdvertisingParams{},
		mVendorId:             nil,
		mProductId:            nil,
		mDeviceType:           nil,
		mPairingHint:          0,
		mPairingInstr:         "",
		mDeviceName:           "",
		mMode:                 0,
		mCommissioningMode:    0,
		mPeerId:               nil,
		mShortDiscriminator:   0,
		mLongDiscriminator:    0,
		mRotatingId:           "",
	}
}

func (c *CommissionAdvertisingParameters) SetCommissioningMode(mode int) {
	c.mCommissioningMode = mode
}

func (c *CommissionAdvertisingParameters) GetCommissioningMode() int {
	return c.mCommissioningMode
}

func (c *CommissionAdvertisingParameters) SetCommissionAdvertiseMode(mode uint16) {
	c.mMode = mode
}

func (c *CommissionAdvertisingParameters) GetCommissionAdvertiseMode() uint16 {
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

func (c *CommissionAdvertisingParameters) SetPairingHint(value uint16) *CommissionAdvertisingParameters {
	c.mPairingHint = value
	return c
}

func (c *CommissionAdvertisingParameters) SetPairingInstruction(ist string) {
	c.mPairingInstr = ist
}

func (c *CommissionAdvertisingParameters) SetMRPConfig(config *messageing.ReliableMessageProtocolConfig) {
	c.mMRPConfig = *config
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

func (c *CommissionAdvertisingParameters) GetPairingHint() uint16 {
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
	c.mLocalMRPConfig = *config
	return c
}
