package main

type Product struct {
	ID string `json:"id"`
	Name string `json:"name"`
	PriceCents int64 `json:"price_cents"`
}
