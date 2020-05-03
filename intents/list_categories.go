package intents

import (
	"fmt"
	"mania/dialogflow"
	"strings"
)

// ListCategoriesHandler handles list_categories intent
func (d *Dispatcher) ListCategoriesHandler(req dialogflow.Request) (dialogflow.Response, error) {
	sess := d.sessions.GetSession(req.Session)
	cats := d.cache.GetCategoriesPage(sess.CurrentPage, d.pageSize)

	if sess.CurrentPage > 0 && len(cats) == 0 {
		d.sessions.ResetPage(req.Session)
		return d.ListCategoriesHandler(req)
	}
	catNames := make([]string, len(cats))
	for i, cat := range cats {
		catNames[i] = cat.Name
	}
	catList := strings.Join(catNames, ", ")
	text := fmt.Sprintf(
		`Вот некоторые из категорий: %s.
Назовите категорию, чтобы посмотреть товары в ней.
Скажите дальше, чтобы вывести ещё категории.`,
		catList)

	// skip details for next pages
	if sess.CurrentPage > 0 {
		text = catList
	}
	resp := dialogflow.GenerateResponse(true, text)

	return resp, nil
}

// ListCategoriesNextHandler handles list_categories_next intent
func (d *Dispatcher) ListCategoriesNextHandler(req dialogflow.Request) (dialogflow.Response, error) {
	d.sessions.NextPage(req.Session)
	return d.ListCategoriesHandler(req)
}
