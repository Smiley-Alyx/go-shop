package main

import "sync"

var (
	// ordersMu защищает доступ к in-memory данным заказов.
	ordersMu    sync.Mutex
	// orders — список заказов в памяти процесса.
	orders      []Order
	// nextOrderID — следующий ID для нового заказа.
	nextOrderID int
)

// storeOrdersInit инициализирует in-memory хранилище (один раз).
func storeOrdersInit() {
	ordersMu.Lock()
	defer ordersMu.Unlock()

	if nextOrderID != 0 {
		return
	}

	orders = nil
	nextOrderID = 1
}

// storeOrderAdd добавляет заказ и назначает ему ID.
func storeOrderAdd(o Order) Order {
	ordersMu.Lock()
	defer ordersMu.Unlock()

	o.ID = nextOrderID
	nextOrderID++

	orders = append(orders, o)
	return o
}

// storeOrderGetByID ищет заказ по ID.
func storeOrderGetByID(id int) (Order, int) {
	ordersMu.Lock()
	defer ordersMu.Unlock()

	for i := 0; i < len(orders); i++ {
		if orders[i].ID == id {
			return orders[i], 1
		}
	}

	return Order{}, 0
}

// storeOrderUpdateStatus меняет статус заказа.
func storeOrderUpdateStatus(id int, status OrderStatus) (Order, int) {
	ordersMu.Lock()
	defer ordersMu.Unlock()

	for i := 0; i < len(orders); i++ {
		if orders[i].ID != id {
			continue
		}

		if orders[i].Status != OrderStatusNew {
			return Order{}, 0
		}
		if status != OrderStatusPaid && status != OrderStatusCancelled {
			return Order{}, 0
		}

		orders[i].Status = status
		return orders[i], 1
	}

	return Order{}, 0
}

// storeOrderList возвращает копию списка заказов.
func storeOrderList() []Order {
	ordersMu.Lock()
	defer ordersMu.Unlock()

	out := make([]Order, len(orders))
	copy(out, orders)
	return out
}
