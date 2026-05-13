package storage

// Session holds auth info persisted across app launches.
type Session struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

// SessionStore persists Session information.
type SessionStore interface {
	Load() (Session, error)
	Save(Session) error
	Clear() error
}

// NewSessionStore returns the platform-specific session store.
func NewSessionStore() SessionStore {
	return newSessionStore()
}
