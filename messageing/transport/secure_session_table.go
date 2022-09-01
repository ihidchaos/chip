package transport

import "math/rand"

type SecureSessionTable struct {
	mNextSessionId uint16
	mEntries       []*SecureSession
}

func NewSecureSessionTable() *SecureSessionTable {
	return &SecureSessionTable{
		mNextSessionId: 0,
		mEntries:       make([]*SecureSession, 0),
	}
}

func (t *SecureSessionTable) Init() {
	t.mNextSessionId = uint16(rand.Uint32())
}

// FindSecureSessionByLocalKey 遍历所有的SecureSession,如果SessionId相同,则取出来
func (t *SecureSessionTable) FindSecureSessionByLocalKey(id uint16) *SecureSession {
	for _, e := range t.mEntries {
		if e.GetLocalSessionId() == id {
			return e
		}
	}
	return nil
}

func (t *SecureSessionTable) CreateNewSecureSession(sessionType TypeSecureSession, sessionId uint16) *SecureSession {
	secureSession := NewSecureSession(t, sessionType, sessionId)
	t.mEntries = append(t.mEntries, secureSession)
	return secureSession
}
