package store

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetCategoriesPage(t *testing.T) {
	ctx := context.Background()
	c, err := NewCache(ctx)
	if err != nil {
		t.Fatalf("unexpected error in NewCache: %v", err)
	}

	cats1 := c.GetCategoriesPage(0, 2)
	if len(cats1) != 2 {
		t.Fatalf("unexpected len=%d from GetCategoriesPage (expected 2)", len(cats1))
	}

	cats2 := c.GetCategoriesPage(1, 1)
	if len(cats2) != 1 {
		t.Fatalf("unexpected len=%d from GetCategoriesPage (expected 1)", len(cats2))
	}

	if diff := cmp.Diff(cats1[1], cats2[0]); diff != "" {
		t.Errorf("cats do differ:\n%s", diff)
	}

	//spew.Dump(cats1)
}

func TestGetItemsPage(t *testing.T) {
	ctx := context.Background()
	c, err := NewCache(ctx)
	if err != nil {
		t.Fatalf("unexpected error in NewCache: %v", err)
	}

	items, err := c.GetItemsPage("Блины", 0, 3)
	if err != nil {
		t.Fatalf("unexpected error in GetItemsPage: %v", err)
	}

	//spew.Dump(items)

	if len(items) != 3 {
		t.Fatalf("unexpected len=%d from GetItemsPage (expected 3)", len(items))
	}
}

func TestGetItem(t *testing.T) {
	ctx := context.Background()
	c, err := NewCache(ctx)
	if err != nil {
		t.Fatalf("unexpected error in NewCache: %v", err)
	}

	item, err := c.GetItem("Блинчики")
	if err != nil {
		t.Fatalf("unexpected error in GetItem: %v", err)
	}

	//spew.Dump(item)

	if item.Name != "Блинчики" {
		t.Fatalf("unexpected item.Name=%s from GetItem (expected Блинчики)", item.Name)
	}
}
