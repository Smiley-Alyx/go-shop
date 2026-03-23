package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

const catalogBaseURL = "http://localhost:8081"

type catalogProduct struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

// fetchProductByID делает HTTP запрос в catalog и возвращает товар.
func fetchProductByID(id int) (catalogProduct, int) {
	client := &http.Client{Timeout: 2 * time.Second}

	req, err := http.NewRequest("GET", catalogBaseURL+"/products/"+strconv.Itoa(id), nil)
	if err != nil {
		return catalogProduct{}, 0
	}

	resp, err := client.Do(req)
	if err != nil {
		return catalogProduct{}, 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return catalogProduct{}, 0
	}

	var p catalogProduct
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&p)
	if err != nil {
		return catalogProduct{}, 0
	}

	return p, 1
}
