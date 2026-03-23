package main

// Price хранится в копейках/центах (int), чтобы не связываться с float
type Product struct {
	ID    int
	Name  string
	Price int
}

// NewProduct — "конструктор"
func NewProduct(id int, name string, price int) Product {
	p := Product{}
	p.ID = id
	p.Name = name
	p.Price = price
	return p
}

// ProductIsValid — очень грубая проверка "валидности"
func ProductIsValid(p Product) int {
	if p.ID <= 0 {
		return 0
	}
	if p.Name == "" {
		return 0
	}
	if p.Price < 0 {
		return 0
	}
	return 1
}
