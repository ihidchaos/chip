package parameters

import (
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/pkg"
	"net"
)

type BaseAdvertisingParams struct {
	mPort           uint16
	mMac            net.HardwareAddr
	mEnableIPv4     bool
	mInterfaceId    net.Interface
	mMRPConfig      *messageing.ReliableMessageProtocolConfig
	mTcpSupported   *bool
	mLocalMRPConfig *messageing.ReliableMessageProtocolConfig
}

func (b BaseAdvertisingParams) Init() *BaseAdvertisingParams {
	b.mMRPConfig = messageing.ReliableMessageProtocolConfig{}.Init()
	b.mLocalMRPConfig = messageing.ReliableMessageProtocolConfig{}.Init()
	return &b
}

func (b *BaseAdvertisingParams) GetLocalMRPConfig() *messageing.ReliableMessageProtocolConfig {
	return b.mLocalMRPConfig
}

func (b *BaseAdvertisingParams) SetLocalMRPConfig(config *messageing.ReliableMessageProtocolConfig) {
	b.mLocalMRPConfig = config
}

func (b *BaseAdvertisingParams) IsIPv4Enabled() bool {
	return b.mEnableIPv4
}

func (b *BaseAdvertisingParams) SetPort(port uint16) {
	b.mPort = port
}

func (b *BaseAdvertisingParams) SetInterfaceId(id net.Interface) {
	b.mInterfaceId = id
}

func (b *BaseAdvertisingParams) GetPort() uint16 {
	return b.mPort
}

func (b *BaseAdvertisingParams) SetMaC(mac net.HardwareAddr) {
	b.mMac = mac
}

func (b *BaseAdvertisingParams) GetMac() net.HardwareAddr {
	if b.mMac == nil || len(b.mMac) == 0 {
		b.mMac = net.HardwareAddr(pkg.Mac48Address(pkg.RandHex()))
	}
	return b.mMac
}

func (b *BaseAdvertisingParams) GetUUID() string {
	if b.mMac == nil || len(b.mMac) == 0 {
		b.mMac = net.HardwareAddr(pkg.Mac48Address(pkg.RandHex()))
	}
	return b.mMac.String()
}

func (b *BaseAdvertisingParams) EnableIpV4(enable bool) {
	b.mEnableIPv4 = enable
}

func (b *BaseAdvertisingParams) GetTcpSupported() *bool {
	return b.mTcpSupported
}

func (b *BaseAdvertisingParams) SetTcpSupported(s bool) {
	b.mTcpSupported = &s
}
