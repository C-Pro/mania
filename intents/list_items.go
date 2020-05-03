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
или 'далее', чтобы вы вывести ещё.`,
		itemList)

	// skip details for next pages
	if sess.CurrentPage > 0 {
		text = itemList
	}
	resp := dialogflow.GenerateResponse(true, text)

	return resp, nil
}

// ListCategoryItemsNextHandler handles list_category_items_next intent
func (d *Dispatcher) ListCategoryItemsNextHandler(req dialogflow.Request) (dialogflow.Response, error) {
	d.sessions.NextPage(req.Session)
	return d.ListCategoryItemsHandler(req)
}

// GetItemHandler handles get_category_item intent
func (d *Dispatcher) GetItemHandler(req dialogflow.Request) (dialogflow.Response, error) {
	itemName, ok := req.QueryResult.Parameters["item"].(string)
	if !ok {
		return dialogflow.GenerateResponse(true, "Не могу распознать блюдо"), nil
	}

	item, err := d.cache.GetItem(itemName)
	if err != nil {
		return dialogflow.GenerateResponse(false, "Не удалось получить информацию о блюде"), err
	}

	text := fmt.Sprintf("%s\n%s\nЦена: %5.2fр.",
		item.Description,
		item.Composition,
		item.Price)
	resp := dialogflow.GenerateResponse(true, text)

	return resp, nil
}
