package intents

import (
	"errors"

	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

var ENoIntent = errors.New("no intent handler found")

type IntentHandler func(dialogflow.WebhookRequest) (dialogflow.WebhookResponse, error)

type IntentName string

const (
	ListCategories    IntentName = "list_categories"
	ListCategoryItems IntentName = "get_category_items"
)

var intentMap = map[IntentName]IntentHandler{
	ListCategories:    ListCategoriesHandler,
	ListCategoryItems: ListCategoryItemsHandler,
}

func GetHandler(displayName string) (IntentHandler, error) {
	h, ok := intentMap[IntentName(displayName)]
	if !ok {
		return nil, ENoIntent
	}
	return h, nil
}
