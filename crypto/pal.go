package crypto

import (
	"crypto/aes"
	"crypto/rand"
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

func AesCcmEncrypt() {
	cip, _ := aes.NewCipher([]byte{})

}
