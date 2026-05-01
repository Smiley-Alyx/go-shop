package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newOrderMuxForTest() *http.ServeMux {
	a := app{version: "test"}
	resetOrdersForTest()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", a.handleOrdersCreate)
	mux.HandleFunc("GET /orders/{id}", a.handleOrdersGet)
	mux.HandleFunc("GET /orders/{id}/status", a.handleOrdersGetStatus)
	return mux
}

func resetOrdersForTest() {
	ordersMu.Lock()
	defer ordersMu.Unlock()

	orders = nil
	nextOrderID = 1
}

func TestOrdersCreateGetStatus(t *testing.T) {
	stubCatalog := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path == "/products/1" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":1,"name":"tea","price":100}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"not found"}`))
	}))
	defer stubCatalog.Close()

	oldURL := catalogBaseURL
	catalogBaseURL = stubCatalog.URL
	defer func() { catalogBaseURL = oldURL }()

	mux := newOrderMuxForTest()

	body := []byte(`{"items":[{"product_id":1,"qty":2}]}`)
	r := httptest.NewRequest("POST", "/orders", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%q", w.Code, w.Body.String())
	}

	var created Order
	err := json.Unmarshal(w.Body.Bytes(), &created)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if created.ID != 1 {
		t.Fatalf("id=%d", created.ID)
	}
	if created.Total != 200 {
		t.Fatalf("total=%d", created.Total)
	}
	if created.Status != OrderStatusNew {
		t.Fatalf("status=%q", created.Status)
	}

	w2 := httptest.NewRecorder()
	r2 := httptest.NewRequest("GET", "/orders/1", nil)
	mux.ServeHTTP(w2, r2)
	if w2.Code != http.StatusOK {
		t.Fatalf("status=%d body=%q", w2.Code, w2.Body.String())
	}

	var got Order
	err = json.Unmarshal(w2.Body.Bytes(), &got)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("id=%d", got.ID)
	}
	if got.Total != 200 {
		t.Fatalf("total=%d", got.Total)
	}

	w3 := httptest.NewRecorder()
	r3 := httptest.NewRequest("GET", "/orders/1/status", nil)
	mux.ServeHTTP(w3, r3)
	if w3.Code != http.StatusOK {
		t.Fatalf("status=%d body=%q", w3.Code, w3.Body.String())
	}

	var st statusResponse
	err = json.Unmarshal(w3.Body.Bytes(), &st)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if st.Status != OrderStatusNew {
		t.Fatalf("status=%q", st.Status)
	}
}

func TestOrdersSetStatus(t *testing.T) {
	stubCatalog := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path == "/products/1" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":1,"name":"tea","price":100}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"not found"}`))
	}))
	defer stubCatalog.Close()

	oldURL := catalogBaseURL
	catalogBaseURL = stubCatalog.URL
	defer func() { catalogBaseURL = oldURL }()

	a := app{version: "test"}
	resetOrdersForTest()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", a.handleOrdersCreate)
	mux.HandleFunc("GET /orders/{id}", a.handleOrdersGet)
	mux.HandleFunc("GET /orders/{id}/status", a.handleOrdersGetStatus)
	mux.HandleFunc("POST /orders/{id}/status", a.handleOrdersSetStatus)

	body := []byte(`{"items":[{"product_id":1,"qty":2}]}`)
	r := httptest.NewRequest("POST", "/orders", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%q", w.Code, w.Body.String())
	}

	body2 := []byte(`{"status":"paid"}`)
	r2 := httptest.NewRequest("POST", "/orders/1/status", bytes.NewReader(body2))
	r2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, r2)

	if w2.Code != http.StatusOK {
		t.Fatalf("set status=%d body=%q", w2.Code, w2.Body.String())
	}

	var updated Order
	err := json.Unmarshal(w2.Body.Bytes(), &updated)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if updated.Status != OrderStatusPaid {
		t.Fatalf("status=%q", updated.Status)
	}

	w3 := httptest.NewRecorder()
	r3 := httptest.NewRequest("GET", "/orders/1/status", nil)
	mux.ServeHTTP(w3, r3)
	if w3.Code != http.StatusOK {
		t.Fatalf("get status=%d body=%q", w3.Code, w3.Body.String())
	}

	var st statusResponse
	err = json.Unmarshal(w3.Body.Bytes(), &st)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if st.Status != OrderStatusPaid {
		t.Fatalf("status=%q", st.Status)
	}

	body4 := []byte(`{"status":"cancelled"}`)
	r4 := httptest.NewRequest("POST", "/orders/1/status", bytes.NewReader(body4))
	r4.Header.Set("Content-Type", "application/json")
	w4 := httptest.NewRecorder()
	mux.ServeHTTP(w4, r4)
	if w4.Code != http.StatusOK {
		t.Fatalf("set cancelled after paid status=%d body=%q", w4.Code, w4.Body.String())
	}

	var cancelled Order
	err = json.Unmarshal(w4.Body.Bytes(), &cancelled)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if cancelled.Status != OrderStatusCancelled {
		t.Fatalf("status=%q", cancelled.Status)
	}

	body6 := []byte(`{"status":"shipped"}`)
	r6 := httptest.NewRequest("POST", "/orders/1/status", bytes.NewReader(body6))
	r6.Header.Set("Content-Type", "application/json")
	w6 := httptest.NewRecorder()
	mux.ServeHTTP(w6, r6)
	if w6.Code != http.StatusBadRequest {
		t.Fatalf("set shipped after cancelled status=%d body=%q", w6.Code, w6.Body.String())
	}

	body5 := []byte(`{"status":"paid"}`)
	r5 := httptest.NewRequest("POST", "/orders/999/status", bytes.NewReader(body5))
	r5.Header.Set("Content-Type", "application/json")
	w5 := httptest.NewRecorder()
	mux.ServeHTTP(w5, r5)
	if w5.Code != http.StatusNotFound {
		t.Fatalf("set status not found status=%d body=%q", w5.Code, w5.Body.String())
	}
}

func TestOrdersSetStatusChainDelivered(t *testing.T) {
	stubCatalog := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path == "/products/1" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":1,"name":"tea","price":100}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"not found"}`))
	}))
	defer stubCatalog.Close()

	oldURL := catalogBaseURL
	catalogBaseURL = stubCatalog.URL
	defer func() { catalogBaseURL = oldURL }()

	a := app{version: "test"}
	resetOrdersForTest()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", a.handleOrdersCreate)
	mux.HandleFunc("POST /orders/{id}/status", a.handleOrdersSetStatus)

	body := []byte(`{"items":[{"product_id":1,"qty":2}]}`)
	r := httptest.NewRequest("POST", "/orders", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%q", w.Code, w.Body.String())
	}

	body2 := []byte(`{"status":"paid"}`)
	r2 := httptest.NewRequest("POST", "/orders/1/status", bytes.NewReader(body2))
	r2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, r2)
	if w2.Code != http.StatusOK {
		t.Fatalf("set paid status=%d body=%q", w2.Code, w2.Body.String())
	}

	body3 := []byte(`{"status":"shipped"}`)
	r3 := httptest.NewRequest("POST", "/orders/1/status", bytes.NewReader(body3))
	r3.Header.Set("Content-Type", "application/json")
	w3 := httptest.NewRecorder()
	mux.ServeHTTP(w3, r3)
	if w3.Code != http.StatusOK {
		t.Fatalf("set shipped status=%d body=%q", w3.Code, w3.Body.String())
	}

	body4 := []byte(`{"status":"delivered"}`)
	r4 := httptest.NewRequest("POST", "/orders/1/status", bytes.NewReader(body4))
	r4.Header.Set("Content-Type", "application/json")
	w4 := httptest.NewRecorder()
	mux.ServeHTTP(w4, r4)
	if w4.Code != http.StatusOK {
		t.Fatalf("set delivered status=%d body=%q", w4.Code, w4.Body.String())
	}

	body5 := []byte(`{"status":"cancelled"}`)
	r5 := httptest.NewRequest("POST", "/orders/1/status", bytes.NewReader(body5))
	r5.Header.Set("Content-Type", "application/json")
	w5 := httptest.NewRecorder()
	mux.ServeHTTP(w5, r5)
	if w5.Code != http.StatusBadRequest {
		t.Fatalf("set cancelled after delivered status=%d body=%q", w5.Code, w5.Body.String())
	}
}
