package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type createOrderRequest struct {
	Items []OrderItem `json:"items"`
}

type apiError struct {
	Error string `json:"error"`
}

type statusResponse struct {
	Status OrderStatus `json:"status"`
}

type setStatusRequest struct {
	Status OrderStatus `json:"status"`
}

func (a app) handleOrdersCreate(w http.ResponseWriter, r *http.Request) {
	var req createOrderRequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad json"})
		return
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad json"})
		return
	}

	if len(req.Items) == 0 {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "empty items"})
		return
	}

	o := NewOrder(0)
	o.Items = req.Items
	o.Total = 0

	for i := 0; i < len(o.Items); i++ {
		if o.Items[i].ProductID <= 0 {
			writeJSON(w, http.StatusBadRequest, apiError{Error: "bad product_id"})
			return
		}
		if o.Items[i].Qty <= 0 {
			writeJSON(w, http.StatusBadRequest, apiError{Error: "bad qty"})
			return
		}

		p, ok := fetchProductByID(o.Items[i].ProductID)
		if ok == 0 {
			writeJSON(w, http.StatusBadRequest, apiError{Error: "product not found"})
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
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad id"})
		return
	}

	o, ok := storeOrderGetByID(id)
	if ok == 0 {
		writeJSON(w, http.StatusNotFound, apiError{Error: "not found"})
		return
	}

	writeJSON(w, http.StatusOK, o)
}

func (a app) handleOrdersGetStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad id"})
		return
	}

	o, ok := storeOrderGetByID(id)
	if ok == 0 {
		writeJSON(w, http.StatusNotFound, apiError{Error: "not found"})
		return
	}

	writeJSON(w, http.StatusOK, statusResponse{Status: o.Status})
}

func (a app) handleOrdersSetStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad id"})
		return
	}

	var req setStatusRequest
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad json"})
		return
	}

	o, ok := storeOrderUpdateStatus(id, req.Status)
	if ok == 0 {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "cannot update status"})
		return
	}

	writeJSON(w, http.StatusOK, o)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	_ = enc.Encode(v)
}
