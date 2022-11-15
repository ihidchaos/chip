package crypto

import (
	"crypto/sha256"
	"golang.org/x/crypto/hkdf"
	"io"
)

func HKDFSha256(mSharedSecret, salt, kKDFSR2Info []byte) []byte {
	hk := hkdf.New(sha256.New, mSharedSecret, salt, kKDFSR2Info)
	data := make([]byte, SymmetricKeyLengthBytes)
	_, err := io.ReadFull(hk, data)
	if err != nil {
		return nil
	}
	return data
}
