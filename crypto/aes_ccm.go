package crypto

import (
	"crypto/aes"
	"github.com/qwerty-iot/dtls/v2"
)

// AES128CCMEncrypt 使用输入加密明文
// 输出的密文长度 =  len(plainText) + tagSize
func AES128CCMEncrypt(plainText, key, nonce, aad []byte, tagSize int) (cipherTag []byte, err error) {
	cipherBook, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	aesCCM, err := dtls.NewCCM(cipherBook, tagSize, len(nonce))
	cipherTag = aesCCM.Seal(nil, nonce, plainText, aad)
	return cipherTag, err
}

func AES128CCMDecrypt(cipherText, key, nonce, addData []byte, tagSize int) (plainText []byte, err error) {
	cipherBook, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	aesCCM, err := dtls.NewCCM(cipherBook, tagSize, len(nonce))
	if err != nil {
		return
	}
	plainText, err = aesCCM.Open(nil, nonce, cipherText, addData)
	return
}
