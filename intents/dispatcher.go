package intents

import (
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
	ListCategories    IntentName = "list_categories"
	ListCategoryItems IntentName = "get_category_items"
)

// Store provides functions to access menu data
type Store interface {
	GetCategoriesPage(pageNum, pageSize int) []*store.Category
	GetItemsPage(categoryName string, pageNum, pageSize int) ([]*store.Item, error)
}

// Dispatcher provides handlers for intents
type Dispatcher struct {
	cache     Store
	intentMap map[IntentName]IntentHandler
}

// NewDispatcher returns new *Dispatcher instance
func NewDispatcher(store Store) *Dispatcher {
	d := Dispatcher{
		cache: store,
	}
	d.intentMap = map[IntentName]IntentHandler{
		ListCategories:    d.ListCategoriesHandler,
		ListCategoryItems: d.ListCategoryItemsHandler,
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
