package dnssd

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	"net"
)

type BaseAdvertisingParams struct {
	mPort           uint16
	mMac            string
	mEnableIPv4     bool
	mInterfaceId    net.Interface
	mMRPConfig      messageing.ReliableMessageProtocolConfig
	mTcpSupported   *bool
	mLocalMRPConfig messageing.ReliableMessageProtocolConfig
}

func NewBaseAdvertisingParams() BaseAdvertisingParams {
	return BaseAdvertisingParams{}
}

func (b BaseAdvertisingParams) Init() BaseAdvertisingParams {
	b.mMRPConfig = messageing.ReliableMessageProtocolConfig{}
	return b
}

func (b *BaseAdvertisingParams) LocalMRPConfig() *messageing.ReliableMessageProtocolConfig {
	return &b.mLocalMRPConfig
}

func (b *BaseAdvertisingParams) SetLocalMRPConfig(config *messageing.ReliableMessageProtocolConfig) {
	b.mLocalMRPConfig = *config
}

func (b *BaseAdvertisingParams) IsIPv4Enabled() bool {
	return b.mEnableIPv4
}

func (b *BaseAdvertisingParams) SetPort(port uint16) *BaseAdvertisingParams {
	b.mPort = port
	return b
}

func (b *BaseAdvertisingParams) SetInterfaceId(id net.Interface) {
	b.mInterfaceId = id
}

func (b *BaseAdvertisingParams) GetPort() uint16 {
	return b.mPort
}

func (b *BaseAdvertisingParams) SetMaC(mac string) {
	b.mMac = mac
}

func (b *BaseAdvertisingParams) GetMac() (string, error) {
	return b.mMac, nil
}

func (b *BaseAdvertisingParams) GetUUID() string {
	return b.mMac
}

func (b *BaseAdvertisingParams) EnableIpV4(enable bool) {
	b.mEnableIPv4 = enable
}

func (b *BaseAdvertisingParams) GetTcpSupported() (bool, error) {
	if b.mTcpSupported == nil {
		return false, lib.IncorrectState
	}
	return *b.mTcpSupported, nil
}

func (b *BaseAdvertisingParams) SetTcpSupported(i int8) {
	var value = true
	if i == 0 {
		value = false
	}
	b.mTcpSupported = &value
}
