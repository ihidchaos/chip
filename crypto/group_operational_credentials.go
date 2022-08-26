package crypto

import "time"

type GroupOperationalCredentials struct {
	StateTime     time.Time
	Hash          uint16
	EncryptionKey []byte
	PrivacyKey    []byte
}
