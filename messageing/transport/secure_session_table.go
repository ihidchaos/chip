package transport

import (
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/messageing/transport/session"
	"math/rand"
)

type SecureSessionTable struct {
	mNextSessionId uint16
	mEntries       [config.SecureSessionPoolSize]*session.Secure
}

func NewSecureSessionTable() *SecureSessionTable {
	return &SecureSessionTable{
		mNextSessionId: 0,
		mEntries:       [config.SecureSessionPoolSize]*session.Secure{},
	}
}

func (t *SecureSessionTable) Init() {
	t.mNextSessionId = uint16(rand.Uint32())
}

func (t *SecureSessionTable) ReleaseSession(session *session.Secure) {
	for i, s := range t.mEntries {
		if s == session {
			t.mEntries[i] = nil
		}
	}
}

// FindSecureSessionByLocalKey 遍历所有的SecureSession,如果SessionId相同,则取出来
func (t *SecureSessionTable) FindSecureSessionByLocalKey(id uint16) *SessionHandle {
	for i, e := range t.mEntries {
		if e.LocalSessionId() == id {
			return NewSessionHandle(t.mEntries[i])
		}
	}
	return nil
}

func (t *SecureSessionTable) CreateSecureSession(sessionType session.SecureType, sessionId uint16) *session.Secure {
	for i, s := range t.mEntries {
		if s == nil {
			t.mEntries[i] = session.NewSecure(t, sessionType, sessionId)
			return t.mEntries[i]
		}
	}
	return nil
}
