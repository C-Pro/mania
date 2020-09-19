package store

import (
	"context"
	"log"
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
	CurrentPage int
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
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
		}

		ss.cleanupExpiredSessions()
	}
}

// NewSession creates new session in the store
func (ss *Sessions) NewSession(id string) {
	ss.mux.Lock()
	defer ss.mux.Unlock()

	log.Printf("L0G: Creating new session %s", id)

	ss.sessions[id] = newSession()
}

// NextPage increments current page for paging operations in the session
func (ss *Sessions) NextPage(id string) {
	ss.mux.Lock()
	defer ss.mux.Unlock()

	s, ok := ss.sessions[id]
	if !ok {
		s = newSession()
	}

	s.CurrentPage++
	ss.sessions[id] = s
}

// ResetPage sets current page for paging operations in the session to zero
func (ss *Sessions) ResetPage(id string) {
	ss.mux.Lock()
	defer ss.mux.Unlock()

	s, ok := ss.sessions[id]
	if ok {
		s.CurrentPage = 0
		ss.sessions[id] = s
	}
}

// AddPosition adds position to user's cart
func (ss *Sessions) AddPosition(id string, pos Position) {
	ss.mux.Lock()
	defer ss.mux.Unlock()

	s, ok := ss.sessions[id]
	if !ok {
		s = newSession()
	}

	log.Printf("L0G: Adding position %v to session %s", pos, id)

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

	log.Printf("L0G: Returning session %s", id)

	s := ss.sessions[id]
	if s != nil {

		log.Printf("L0G: Session %s: %v", id, *s)
		return *s
	}

	log.Printf("L0G: Session %s is empty", id)
	ss.sessions[id] = newSession()

	return *ss.sessions[id]
}
