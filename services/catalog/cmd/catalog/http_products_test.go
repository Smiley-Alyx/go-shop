package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newCatalogMuxForTest() *http.ServeMux {
	a := app{version: "test"}
	resetProductsForTest()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /products", a.handleProductsList)
	mux.HandleFunc("GET /products/{id}", a.handleProductsGet)
	mux.HandleFunc("POST /products", a.handleProductsCreate)
	return mux
}

func resetProductsForTest() {
	productsMu.Lock()
	defer productsMu.Unlock()

	products = nil
	nextProductID = 1
}

func TestProductsCreateAndGet(t *testing.T) {
	mux := newCatalogMuxForTest()

	body := []byte(`{"name":"tea","price":120}`)
	r := httptest.NewRequest("POST", "/products", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%q", w.Code, w.Body.String())
	}

	var created Product
	err := json.Unmarshal(w.Body.Bytes(), &created)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if created.ID != 1 {
		t.Fatalf("id=%d", created.ID)
	}

	r2 := httptest.NewRequest("GET", "/products/1", nil)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, r2)

	if w2.Code != http.StatusOK {
		t.Fatalf("status=%d body=%q", w2.Code, w2.Body.String())
	}

	var got Product
	err = json.Unmarshal(w2.Body.Bytes(), &got)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.ID != 1 {
		t.Fatalf("id=%d", got.ID)
	}
	if got.Name != "tea" {
		t.Fatalf("name=%q", got.Name)
	}
	if got.Price != 120 {
		t.Fatalf("price=%d", got.Price)
	}
}

func TestProductsList(t *testing.T) {
	mux := newCatalogMuxForTest()

	body1 := []byte(`{"name":"tea","price":120}`)
	r1 := httptest.NewRequest("POST", "/products", bytes.NewReader(body1))
	r1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	mux.ServeHTTP(w1, r1)
	if w1.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%q", w1.Code, w1.Body.String())
	}

	body2 := []byte(`{"name":"coffee","price":200}`)
	r2 := httptest.NewRequest("POST", "/products", bytes.NewReader(body2))
	r2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, r2)
	if w2.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%q", w2.Code, w2.Body.String())
	}

	r3 := httptest.NewRequest("GET", "/products", nil)
	w3 := httptest.NewRecorder()
	mux.ServeHTTP(w3, r3)

	if w3.Code != http.StatusOK {
		t.Fatalf("status=%d body=%q", w3.Code, w3.Body.String())
	}

	var items []Product
	err := json.Unmarshal(w3.Body.Bytes(), &items)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("len=%d", len(items))
	}
	if items[0].ID != 1 {
		t.Fatalf("id0=%d", items[0].ID)
	}
	if items[1].ID != 2 {
		t.Fatalf("id1=%d", items[1].ID)
	}
}
