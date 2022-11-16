package transport

import (
	"github.com/galenliu/chip/messageing/transport/session"
	"math/rand"
)
import "golang.org/x/exp/slices"

type SecureSessionTable struct {
	mNextSessionId uint16
	mEntries       []*session.Secure
}

func NewSecureSessionTable() *SecureSessionTable {
	return &SecureSessionTable{
		mNextSessionId: 0,
		mEntries:       make([]*session.Secure, 0),
	}
}

func (t *SecureSessionTable) Init() {
	t.mNextSessionId = uint16(rand.Uint32())
}

func (t *SecureSessionTable) ReleaseSession(session *session.Secure) {
	if index := slices.Index(t.mEntries, session); index >= 0 {
		t.mEntries = slices.Delete(t.mEntries, index, index+1)
	}
}

// FindSecureSessionByLocalKey 遍历所有的SecureSession,如果SessionId相同,则取出来
func (t *SecureSessionTable) FindSecureSessionByLocalKey(id uint16) *SessionHandle {
	for _, e := range t.mEntries {
		if e.LocalSessionId() == id {
			return NewSessionHandle(e)
		}
	}
	return nil
}

func (t *SecureSessionTable) CreateNewSecureSession(sessionType session.SecureSessionType, sessionId uint16) *session.Secure {
	secureSession := session.NewSecure(t, sessionType, sessionId)
	t.mEntries = append(t.mEntries, secureSession)
	return secureSession
}
