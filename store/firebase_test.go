package store

import (
	"testing"

	"golang.org/x/net/context"
)

func TestNew(t *testing.T) {
	ctx := context.Background()
	_, err := New(ctx)
	if err != nil {
		t.Fatalf("Unexpected error in New: %v", err)
	}
}

func TestGetCategories(t *testing.T) {
	ctx := context.Background()
	s, err := New(ctx)
	if err != nil {
		t.Fatalf("Unexpected error in New: %v", err)
	}

	cats, err := s.GetCategories(ctx)
	if err != nil {
		t.Fatalf("Unexpected error in GetCategories: %v", err)
	}

	if len(cats) == 0 {
		t.Error("expected to have some categories")
	}

	//b, _ := json.Marshal(cats)
	//fmt.Println(string(b))
	//spew.Dump(cats)
}

func TestGetItems(t *testing.T) {
	ctx := context.Background()
	s, err := New(ctx)
	if err != nil {
		t.Fatalf("Unexpected error in New: %v", err)
	}

	items, err := s.GetItems(ctx)
	if err != nil {
		t.Fatalf("Unexpected error in GetItems: %v", err)
	}

	if len(items) == 0 {
		t.Error("expected to have some menu items")
	}

	//b, _ := json.Marshal(cats)
	//fmt.Println(string(b))
	//spew.Dump(items)
}
