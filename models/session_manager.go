package models

import "sync"

type SessionManager struct {
	sessions map[string]*Session
	mutex    *sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
		mutex:    new(sync.RWMutex),
	}
}

func (sm *SessionManager) Get_session(sessionID string) (*Session, bool) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	session, ok := sm.sessions[sessionID]
	return session, ok
}

func (sm *SessionManager) Set_session(session *Session) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.sessions[session.Session_token] = session
}
