package crypto

import (
	"crypto/sha256"
	"golang.org/x/crypto/hkdf"
	"io"
)

func HKDFSha256(mSharedSecret, salt, kKDFSR2Info []byte) []byte {
	buf := hkdf.New(sha256.New, mSharedSecret, salt, kKDFSR2Info)
	data := make([]byte, SymmetricKeyLengthBytes)
	_, err := io.ReadFull(buf, data)
	if err != nil {
		return data
	}
	return nil
}
