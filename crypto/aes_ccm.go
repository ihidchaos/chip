package crypto

import (
	"crypto/aes"
	"github.com/CrimsonAIO/aesccm"
)

func AESCCMEncrypt(plainText, key, nonce, data []byte, tagSize int) (cipherText []byte, err error) {

	cipherBook, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	aesCCM, err := aesccm.NewCCM(cipherBook, len(nonce), tagSize)
	if err != nil {
		return
	}
	cipherText = aesCCM.Seal(nil, nonce, plainText, data)
	return
}

func AESCCMDecrypt(cipherText, key, nonce, data []byte, tagSize int) (plainText []byte, err error) {

	cipherBook, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	aesCCM, err := aesccm.NewCCM(cipherBook, len(nonce), tagSize)
	if err != nil {
		return
	}
	plainText, err = aesCCM.Open(nil, nonce, cipherText, data)
	return
}
