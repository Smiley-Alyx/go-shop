package main

// OrderStatus — статус заказа
type OrderStatus string

const (
	OrderStatusNew       OrderStatus = "new"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// OrderItem — позиция заказа
type OrderItem struct {
	ProductID int
	Qty       int
}

// Order — модель заказа.
// Total — сумма заказа в копейках/центах (int). Пока считаем как хотим, без точности/валют
type Order struct {
	ID     int
	Status OrderStatus
	Items  []OrderItem
	Total  int
}

// NewOrder — "конструктор"
func NewOrder(id int) Order {
	o := Order{}
	o.ID = id
	o.Status = OrderStatusNew
	o.Items = nil
	o.Total = 0
	return o
}

// OrderIsValid - очень грубая проверка
func OrderIsValid(o Order) int {
	if o.ID <= 0 {
		return 0
	}
	if o.Status == "" {
		return 0
	}
	// Items и Total пока не проверяем — всё равно будем переписывать.
	return 1
}
