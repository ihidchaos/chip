package crypto

import (
	"crypto/sha256"
	"testing"
)

func TestHKDFSha256(t *testing.T) {
	hash := sha256.New
	secret := []byte{0x00, 0x01, 0x02, 0x03}
	salt := make([]byte, hash().Size())
	info := []byte("hkdf example")
	data := HKDFSha256(secret, salt, info)
	t.Logf("HKDF: %0X", data)
}
