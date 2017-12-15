package network

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// 会话访问
type SessionAccessor interface {

	// 获取一个连接
	GetSession(int64) Session

	// 遍历连接
	VisitSession(func(Session) bool)

	// 连接数量
	SessionCount() int

	// 关闭所有连接
	CloseAllSession()
}

// 完整功能的会话管理
type SessionManager interface {
	SessionAccessor

	Add(Session) error
	Remove(Session)
}

type SessionManagerImplement struct {
	sessionMap      map[int64]Session
	sessionIDAcc    int64
	sessionMapGuard sync.RWMutex
}

const totalTryCount = 100

var ErrSessionIDOverride = fmt.Errorf("Session ID override")

func (self *SessionManagerImplement) Add(s Session) error {

	self.sessionMapGuard.Lock()
	defer self.sessionMapGuard.Unlock()
	var tryCount int = totalTryCount
	var id int64

	for tryCount > 0 {
		id = atomic.AddInt64(&self.sessionIDAcc, 1)
		if _, ok := self.sessionMap[id]; !ok {
			break
		}
		tryCount--
	}
	if tryCount == 0 {
		return ErrSessionIDOverride
	}
	s.(interface {
		SetID(int64)
	}).SetID(id)
	self.sessionMap[id] = s
	return nil
}
func (self *SessionManagerImplement) Remove(s Session) {
	self.sessionMapGuard.Lock()
	delete(self.sessionMap, s.ID())
	self.sessionMapGuard.Unlock()
}
func (self *SessionManagerImplement) GetSession(id int64) Session {
	self.sessionMapGuard.RLock()
	defer self.sessionMapGuard.RUnlock()
	v, ok := self.sessionMap[id]
	if ok {
		return v
	}
	return nil
}

func (self *SessionManagerImplement) SessionCount() int {
	self.sessionMapGuard.Lock()
	defer self.sessionMapGuard.Unlock()

	return len(self.sessionMap)
}

func (self *SessionManagerImplement) VisitSession(callback func(Session) bool) {
	self.sessionMapGuard.RLock()
	defer self.sessionMapGuard.RUnlock()

	for _, ses := range self.sessionMap {
		if !callback(ses) {
			break
		}
	}
}

func (self *SessionManagerImplement) CloseAllSession() {
	self.VisitSession(func(ses Session) bool {
		ses.Close()
		return true
	})
}
func NewSessionManager() SessionManager {
	return &SessionManagerImplement{
		sessionMap: make(map[int64]Session),
	}
}
