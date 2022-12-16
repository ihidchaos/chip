package transport

import (
	"bytes"
	"encoding/binary"
	"github.com/galenliu/chip"
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/messageing/transport/session"
	"github.com/galenliu/chip/platform/system"
)

type KeyUsage uint8

const (
	kMaxAADLen            = 128
	KSessionRoleInitiator = iota
	KSessionRoleResponder
	kAESCCMNonceLen          uint8    = 13
	kI2RKey                  KeyUsage = 0
	kR2IKey                  KeyUsage = 1
	kAttestationChallengeKey KeyUsage = 2
	kNumCryptoKeys           KeyUsage = 3
)

type CryptoKey [crypto.AESCCM128KeyLength]byte

type CryptoContext struct {
	mKeyContext   crypto.SymmetricKeyContextBase
	mKeyAvailable bool
	mKeys         [kNumCryptoKeys]CryptoKey
	mSessionRole  session.Role
}

func NewCryptoContext(key crypto.SymmetricKeyContextBase) *CryptoContext {
	return &CryptoContext{}
}

func (c *CryptoContext) Decrypt(msg *system.PacketBufferHandle, nonce []byte, header *raw.PacketHeader, mac *raw.MessageAuthenticationCode) {
	_ = mac.Tag()

	//header.Ecode()
}

func (c *CryptoContext) Encrypt(plainText []byte, nonce []byte, header *raw.PacketHeader) (cipherTag []byte, err error) {

	add, err := c.GetAdditionalAuthData(header)
	if err != nil {
		return nil, err
	}
	if c.mKeyContext != nil {
		if cipherTag, err := c.mKeyContext.MessageEncrypt(plainText, nonce, add, header.MICTagLength()); err != nil {
			return nil, err
		} else {
			return cipherTag, nil
		}

	} else {
		if !c.mKeyAvailable {
			return nil, chip.ErrorInvalidUseOfSessionKey
		}
		usage := kR2IKey
		if c.mSessionRole == session.RoleInitiator {
			usage = kI2RKey
		}
		if cipherTag, err := crypto.AES128CCMEncrypt(plainText, c.mKeys[usage][:], nonce, add, header.MICTagLength()); err != nil {
			return nil, err
		} else {
			return cipherTag, nil
		}
	}
}

func (c *CryptoContext) GetAdditionalAuthData(header *raw.PacketHeader) ([]byte, error) {
	encodeBuf := bytes.NewBuffer(nil)
	if err := header.Encode(encodeBuf); err != nil {
		return nil, err
	}
	return encodeBuf.Bytes(), nil
}

// BuildNonce 使用SecFlags,messageCounter,nodeId三个字段生成Nonce(用于AES加解密的初始化向量)，Len == 13
func BuildNonce(secFlags uint8, messageCounter uint32, nodeId lib.NodeId) []byte {
	nonceStorage := make([]byte, kAESCCMNonceLen)
	nonceStorage[0] = secFlags
	binary.LittleEndian.PutUint32(nonceStorage[1:5], messageCounter)
	binary.LittleEndian.PutUint64(nonceStorage[5:12], uint64(nodeId))
	return nonceStorage
}
