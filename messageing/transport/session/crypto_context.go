package session

import (
	"encoding/binary"
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/platform/system"
)

const (
	KSessionRoleInitiator = iota
	KSessionRoleResponder
	kAESCCMNonceLen uint8 = 13
)

type CryptoContext struct {
}

func NewCryptoContext(key *crypto.SymmetricKeyContext) *CryptoContext {
	return &CryptoContext{}
}

func (c *CryptoContext) Decrypt(msg *system.PacketBufferHandle, nonce []byte, header *raw.PacketHeader, mac *raw.MessageAuthenticationCode) {
	_ = mac.Tag()

	//header.Ecode()
}

// BuildNonce 使用SecFlags,messageCounter,nodeId三个字段生成Nonce(用于AES加解密的初始化向量)，Len == 13
func BuildNonce(secFlags uint8, messageCounter uint32, nodeId lib.NodeId) []byte {
	nonceStorage := make([]byte, kAESCCMNonceLen)
	nonceStorage[0] = secFlags
	binary.LittleEndian.PutUint32(nonceStorage[1:5], messageCounter)
	binary.LittleEndian.PutUint64(nonceStorage[5:12], uint64(nodeId))
	return nonceStorage
}

func GetAdditionalAuthData(header *raw.PacketHeader) {

}
