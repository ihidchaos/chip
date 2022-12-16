package credentials

import "github.com/galenliu/chip/crypto"

type GroupDataProviderImpl struct {
}

type GroupKeyContext struct {
	mEncryptionKey []byte
}

func (g GroupKeyContext) KeyHash() uint16 {
	//TODO implement me
	panic("implement me")
}

func (g GroupKeyContext) Release() {
	//TODO implement me
	panic("implement me")
}

func (g GroupKeyContext) MessageEncrypt(plaintext, nonce, addData []byte, tagSize int) (cipherTag []byte, err error) {
	return crypto.AES128CCMEncrypt(plaintext, g.mEncryptionKey, nonce, addData, tagSize)
}
