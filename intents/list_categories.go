package intents

import (
	"mania/dialogflow"
	"strings"
)

// ListCategoriesHandler handles list_categories intent
func (d *Dispatcher) ListCategoriesHandler(req dialogflow.Request) (dialogflow.Response, error) {
	cats := d.cache.GetCategoriesPage(0, 10)
	catNames := make([]string, len(cats))
	for i, cat := range cats {
		catNames[i] = cat.Name
	}
	catList := strings.Join(catNames, ", ")
	resp := dialogflow.GenerateResponse(true, catList)

	return resp, nil
}
