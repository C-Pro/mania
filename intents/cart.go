package intents

import (
	"fmt"
	"mania/dialogflow"
	"mania/store"
	"strconv"
)

// AddToCartHandler handles add_to_cart_context intent
func (d *Dispatcher) AddToCartHandler(req dialogflow.Request) (dialogflow.Response, error) {

	itemName, ok := req.QueryResult.Parameters["item"].(string)
	if !ok {
		return dialogflow.GenerateResponse(true, "Не могу распознать блюдо"), nil
	}

	item, err := d.cache.GetItem(itemName)
	if err != nil {
		return dialogflow.GenerateResponse(false, "Не удалось получить информацию о блюде"), err
	}


	quantity := uint(1)
	numberStr, ok := req.QueryResult.Parameters["number"].(string)
	if ok {
		qint, err := strconv.Atoi(numberStr)
		if err == nil && qint > 0 && qint <= 100 {
			quantity = uint(qint)
		}
	}

	d.sessions.AddPosition(
		req.Session,
		store.Position{
			Item:     *item,
			Quantity: quantity,
		})

	sess := d.sessions.GetSession(req.Session)
	cnt := uint(0)
	amount := 0.0
	for _, pos := range sess.Cart {
		cnt += pos.Quantity
		amount += pos.Item.Price
	}

	text := fmt.Sprintf("В корзине %d товаров на сумму %5.2f рублей", cnt, amount)

	return dialogflow.GenerateResponse(true, text), nil
}
