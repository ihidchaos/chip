package transport

import (
	"encoding/binary"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport/raw"
)

const (
	KSessionRoleInitiator = iota
	KSessionRoleResponder

	kAESCCMNonceLen uint8 = 13
)

type CryptoContext struct {
}

func (c CryptoContext) Decrypt(msg *raw.PacketBuffer, nonce []byte, header *raw.PacketHeader, mac *raw.MessageAuthenticationCode) {
	_ = mac.Tag()

	//header.Ecode()

}

// BuildNonce 使用SecFlags,messageCounter,nodeId三个字段生成Nonce，Len == 13
func BuildNonce(secFlags uint8, messageCounter uint32, nodeId lib.NodeId) ([]byte, error) {
	data := make([]byte, 13)
	data[0] = secFlags
	binary.LittleEndian.PutUint32(data[1:5], messageCounter)
	binary.LittleEndian.PutUint64(data[5:12], uint64(nodeId))
	return data, nil
}

func GetAdditionalAuthData(header *raw.PacketHeader) {

}
