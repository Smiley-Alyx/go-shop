package main

import "sync"

var (
	// productsMu защищает доступ к in-memory данным каталога.
	productsMu   sync.Mutex
	// products — список товаров в памяти процесса.
	products     []Product
	// nextProductID — следующий ID для нового товара.
	nextProductID int
)

// storeProductsInit инициализирует in-memory хранилище (один раз).
func storeProductsInit() {
	productsMu.Lock()
	defer productsMu.Unlock()

	if nextProductID != 0 {
		return
	}

	products = nil
	nextProductID = 1
}

// storeProductAdd добавляет товар и назначает ему ID.
func storeProductAdd(p Product) Product {
	productsMu.Lock()
	defer productsMu.Unlock()

	p.ID = nextProductID
	nextProductID++

	products = append(products, p)
	return p
}

// storeProductGetByID ищет товар по ID.
func storeProductGetByID(id int) (Product, int) {
	productsMu.Lock()
	defer productsMu.Unlock()

	for i := 0; i < len(products); i++ {
		if products[i].ID == id {
			return products[i], 1
		}
	}

	return Product{}, 0
}

// storeProductList возвращает копию списка товаров.
func storeProductList() []Product {
	productsMu.Lock()
	defer productsMu.Unlock()

	out := make([]Product, len(products))
	copy(out, products)
	return out
}
