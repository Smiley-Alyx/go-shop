package main

// OrderStatus — строковый статус заказа.
type OrderStatus string

const (
	OrderStatusNew       OrderStatus = "new"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// OrderItem — позиция заказа.
type OrderItem struct {
	ProductID int
	Qty       int
}

// Order — модель заказа.
// Total — сумма заказа в копейках/центах (int).
type Order struct {
	ID     int
	Status OrderStatus
	Items  []OrderItem
	Total  int
}

// NewOrder создаёт заказ со статусом new.
func NewOrder(id int) Order {
	o := Order{}
	o.ID = id
	o.Status = OrderStatusNew
	o.Items = nil
	o.Total = 0
	return o
}

// OrderIsValid делает простую проверку заполненности обязательных полей.
func OrderIsValid(o Order) int {
	if o.ID <= 0 {
		return 0
	}
	if o.Status == "" {
		return 0
	}
	// Items и Total здесь не проверяем.
	return 1
}
