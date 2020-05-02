package intents

import (
	"fmt"
	"mania/dialogflow"
	"strings"
)

// ListCategoryItemsHandler handles get_category_items intent
func (d *Dispatcher) ListCategoryItemsHandler(req dialogflow.Request) (dialogflow.Response, error) {
	categoryName, ok := req.QueryResult.Parameters["category"].(string)
	if !ok {
		return dialogflow.GenerateResponse(true, "Не могу распознать категорию меню"), nil
	}
	sess := d.sessions.GetSession(req.Session)
	items, err := d.cache.GetItemsPage(categoryName, sess.CurrentPage, d.pageSize)
	if err != nil {
		return dialogflow.GenerateResponse(false, "Не удалось получить содержимое категории"), err
	}

	if len(items) == 0 {
		if sess.CurrentPage > 0 {
			d.sessions.ResetPage(req.Session)
			return d.ListCategoryItemsHandler(req)
		}
		return dialogflow.GenerateResponse(
			true,
			fmt.Sprintf("Нет позиций в категории %s", categoryName),
		), nil
	}

	itemNames := make([]string, len(items))
	for i, item := range items {
		itemNames[i] = item.Name
	}
	itemList := strings.Join(itemNames, ", ")
	text := fmt.Sprintf(
		`%s.
Назовите продукт, чтобы узнать о нём подробней или добавить в корзину,
или скажите 'дальше', чтобы вы вести ещё.`,
		itemList)
	resp := dialogflow.GenerateResponse(true, text)

	return resp, nil
}
