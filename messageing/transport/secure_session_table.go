package transport

import "math/rand"

type SecureSessionTable struct {
	mNextSessionId uint16
	mEntries       []*SecureSession
}

func NewSecureSessionTable() *SecureSessionTable {
	return &SecureSessionTable{}
}

func (t *SecureSessionTable) Init() {
	t.mNextSessionId = uint16(rand.Uint32())
}

func (t *SecureSessionTable) FindSecureSessionByLocalKey(id uint16) *SecureSession {
	for _, e := range t.mEntries {
		if e.GetLocalSessionId() == id {
			return e
		}
	}
	return nil
}
