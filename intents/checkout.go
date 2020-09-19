package intents

import (
	"fmt"
	"mania/dialogflow"
)

// CheckoutHandler handles checkout_intent
func (d *Dispatcher) CheckoutHandler(req dialogflow.Request) (dialogflow.Response, error) {

	phoneNumber, ok := req.QueryResult.Parameters["phonenum"].(string)
	if !ok {
		return dialogflow.GenerateResponse(true, "Укажите номер телефона"), nil
	}

	sess := d.sessions.GetSession(req.Session)
	if len(sess.Cart) == 0 {
		return dialogflow.GenerateResponse(false, "Корзина пуста"), nil
	}

	cnt := uint(0)
	amount := 0.0

	items := ""

	for _, pos := range sess.Cart {
		cnt += pos.Quantity
		amount += pos.Item.Price
		items = fmt.Sprintf("%s%s - %dшт\n", items, pos.Item.Name, pos.Quantity)
	}

	text := fmt.Sprintf("Заказ от %s: %d товаров на сумму %5.2f рублей:\n%s", phoneNumber, cnt, amount, items)
	if err := d.Send(text, phoneNumber); err != nil {
		return dialogflow.GenerateResponse(false, "Ошибка отправки заказа, попробуйте ещё"), err
	}

	return dialogflow.GenerateResponse(true, "Ваш заказ зарегистрирован. Ожидайте звонка. Спасибо!"), nil
}
