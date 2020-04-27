package store

import (
	"context"
	"sync"
	"time"
)

const (
	maxTTL         = time.Minute * 30
	expireInterval = time.Minute * 5
)

// Position holds invoice line for users' cart
type Position struct {
	Item     Item
	Quantity uint
}

// Session holds user conversation context
type Session struct {
	CurrentPage uint
	Cart        map[int]Position
	created     time.Time
}

// newSession returns a new Session instance
func newSession() *Session {
	return &Session{
		created: time.Now().UTC(),
		Cart:    make(map[int]Position),
	}
}

// Sessions stores all active users' conversations' contexts
type Sessions struct {
	mux      sync.RWMutex
	sessions map[string]*Session
}

// NewSessions returns a new Sessions store instance
func NewSessions(ctx context.Context) *Sessions {
	s := Sessions{
		sessions: make(map[string]*Session),
	}

	go s.cleanupLoop(ctx)

	return &s
}

// cleanupExpiredSessions removes expired sessions
func (ss *Sessions) cleanupExpiredSessions() {
	ss.mux.Lock()
	defer ss.mux.Unlock()

	for k := range ss.sessions {
		if time.Now().UTC().Sub(ss.sessions[k].created) < maxTTL {
			delete(ss.sessions, k)
		}
	}
}

// cleanupLoop cleans expired sessions every expireInterval
func (ss *Sessions) cleanupLoop(ctx context.Context) {
	tick := time.NewTicker(expireInterval)

	for {
		ss.cleanupExpiredSessions()

		select {
		case <-ctx.Done():
			break
		case <-tick.C:
		}
	}
}

// NewSession creates new session in the store
func (ss *Sessions) NewSession(id string) {
	ss.mux.Lock()
	defer ss.mux.Unlock()

	ss.sessions[id] = newSession()
}

// SetPage Sets current page for paging operations in the session
func (ss *Sessions) SetPage(id string, currentPage uint) {
	ss.mux.Lock()
	defer ss.mux.Unlock()

	s, ok := ss.sessions[id]
	if !ok {
		s = newSession()
	}

	s.CurrentPage = currentPage
	ss.sessions[id] = s
}

// AddPosition adds position to user's cart
func (ss *Sessions) AddPosition(id string, pos Position) {
	ss.mux.Lock()
	defer ss.mux.Unlock()

	s, ok := ss.sessions[id]
	if !ok {
		s = newSession()
	}

	s.Cart[pos.Item.ID] = pos
}

// RemovePosition removes position to user's cart
func (ss *Sessions) RemovePosition(id string, itemID int) {
	ss.mux.Lock()
	defer ss.mux.Unlock()

	s, ok := ss.sessions[id]
	if ok {
		delete(s.Cart, itemID)
	}
}

// RemoveCart removes user's cart
func (ss *Sessions) RemoveCart(id string, itemID int) {
	ss.mux.Lock()
	defer ss.mux.Unlock()

	s, ok := ss.sessions[id]
	if ok {
		s.Cart = make(map[int]Position)
	}
}

// GetSession returns user's session object
func (ss *Sessions) GetSession(id string) Session {
	ss.mux.RLock()
	defer ss.mux.RUnlock()

	s := ss.sessions[id]
	if s != nil {
		return *s
	}

	return Session{}
}
