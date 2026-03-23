package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type createOrderRequest struct {
	Items []OrderItem `json:"items"`
}

type statusResponse struct {
	Status OrderStatus `json:"status"`
}

func (a app) handleOrdersCreate(w http.ResponseWriter, r *http.Request) {
	var req createOrderRequest

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "bad json"})
		return
	}

	if len(req.Items) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "empty items"})
		return
	}

	o := NewOrder(0)
	o.Items = req.Items
	o.Total = 0

	for i := 0; i < len(o.Items); i++ {
		if o.Items[i].ProductID <= 0 {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "bad product_id"})
			return
		}
		if o.Items[i].Qty <= 0 {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "bad qty"})
			return
		}

		p, ok := fetchProductByID(o.Items[i].ProductID)
		if ok == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "product not found"})
			return
		}

		o.Total += p.Price * o.Items[i].Qty
	}

	o = storeOrderAdd(o)
	writeJSON(w, http.StatusCreated, o)
}

func (a app) handleOrdersGet(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "bad id"})
		return
	}

	o, ok := storeOrderGetByID(id)
	if ok == 0 {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
		return
	}

	writeJSON(w, http.StatusOK, o)
}

func (a app) handleOrdersGetStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "bad id"})
		return
	}

	o, ok := storeOrderGetByID(id)
	if ok == 0 {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
		return
	}

	writeJSON(w, http.StatusOK, statusResponse{Status: o.Status})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	_ = enc.Encode(v)
}
