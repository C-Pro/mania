package store

import (
	"context"
	"testing"
)

func TestNewSessions(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := NewSessions(ctx)
	if s == nil {
		t.Error("Sessions shouldn't be nil")
	}
}

func TestNewSession(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	s := NewSessions(ctx)
	s.NewSession("123")
	if len(s.sessions) != 1 {
		t.Errorf("expected one session, got %d", len(s.sessions))
	}
}

func TestNextPage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	s := NewSessions(ctx)
	s.NewSession("123")
	s.NextPage("123")
	s.NextPage("123")
	s.NextPage("123")
	if len(s.sessions) != 1 {
		t.Errorf("expected one session, got %d", len(s.sessions))
	}
	s2 := s.GetSession("123")
	if s2.CurrentPage != 3 {
		t.Errorf("expected session to have current page = 3, but got %d", s2.CurrentPage)
	}
	s.ResetPage("123")
	s3 := s.GetSession("123")
	if s3.CurrentPage != 0 {
		t.Errorf("expected session to have current page = 0, but got %d", s3.CurrentPage)
	}
}
