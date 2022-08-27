package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

type Spake2pVerifier struct {
}

func (v Spake2pVerifier) Deserialize(verifier []byte) error {
	return nil
}

func (v Spake2pVerifier) Generate(count uint32, span []byte, passcode uint32) error {
	return nil
}

func (v Spake2pVerifier) Serialize() ([]byte, error) {
	return nil, nil
}

type P256ECDSASignature struct {
}

func (s *P256ECDSASignature) Bytes() []byte {
	return nil
}

func SignP256ECDSASignature(plainTex, privateKeyFile []byte) (P256ECDSASignature, error) {
	return P256ECDSASignature{}, nil
}

type P256ECDHDerivedSecret struct {
}

func DRBGBytes(data []byte) error {
	_, err := rand.Read(data)
	if err != nil {
		return err
	}
	return nil
}

// AesCcmEncrypt AES 要求的Key长度为16个字节
func AesCcmEncrypt(plainText, key []byte) (outPut []byte, err error) {
	if len(key) != SymmetricKeyLengthBytes {
		return nil, fmt.Errorf("key length mismatch:%d", len(key))
	}
	return
}

func AesCtrEncryptOrDecrypt(inputText, key, iv []byte) (outPut []byte, err error) {
	//创建一个底层使用AES的密码接口
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	//创建一个使用CTR分组接口, iv的长度等于明文分组的长度
	stream := cipher.NewCTR(block, iv)
	outPut = make([]byte, len(inputText))
	stream.XORKeyStream(outPut, inputText)
	return
}
