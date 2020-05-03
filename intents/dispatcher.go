package intents

import (
	"context"
	"errors"

	"mania/dialogflow"
	"mania/store"
)

// ErrNoIntent is raised when no handler is found for the event
var ErrNoIntent = errors.New("no intent handler found")

// IntentHandler intent handler function type
type IntentHandler func(dialogflow.Request) (dialogflow.Response, error)

// IntentName type for intent name
type IntentName string

// Intent names "enum"
const (
	ListCategories        IntentName = "list_categories"
	ListCategoryItems     IntentName = "get_category_items"
	ListCategoriesNext    IntentName = "list_categories_next"
	ListCategoryItemsNext IntentName = "get_category_items_next"
	GetItem               IntentName = "get_category_item"
)

// Store provides functions to access menu data
type Store interface {
	GetCategoriesPage(pageNum, pageSize int) []*store.Category
	GetItemsPage(categoryName string, pageNum, pageSize int) ([]*store.Item, error)
	GetItem(itemName string) (*store.Item, error)
}

// Dispatcher provides handlers for intents
type Dispatcher struct {
	cache     Store
	sessions  *store.Sessions
	intentMap map[IntentName]IntentHandler
	pageSize  int
}

// NewDispatcher returns new *Dispatcher instance
func NewDispatcher(ctx context.Context, s Store) *Dispatcher {
	d := Dispatcher{
		cache:    s,
		sessions: store.NewSessions(ctx),
		pageSize: 7,
	}
	d.intentMap = map[IntentName]IntentHandler{
		ListCategories:        d.ListCategoriesHandler,
		ListCategoryItems:     d.ListCategoryItemsHandler,
		ListCategoriesNext:    d.ListCategoriesNextHandler,
		ListCategoryItemsNext: d.ListCategoryItemsNextHandler,
		GetItem:               d.GetItemHandler,
	}

	return &d
}

// GetHandler returns a handler for intent webhook
func (d *Dispatcher) GetHandler(intentName string) (IntentHandler, error) {
	h, ok := d.intentMap[IntentName(intentName)]
	if !ok {
		return nil, ErrNoIntent
	}
	return h, nil
}
