package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type createProductRequest struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func (a app) handleProductsList(w http.ResponseWriter, r *http.Request) {
	items := storeProductList()
	writeJSON(w, http.StatusOK, items)
}

func (a app) handleProductsGet(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "bad id"})
		return
	}

	p, ok := storeProductGetByID(id)
	if ok == 0 {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
		return
	}

	writeJSON(w, http.StatusOK, p)
}

func (a app) handleProductsCreate(w http.ResponseWriter, r *http.Request) {
	var req createProductRequest

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "bad json"})
		return
	}

	p := Product{}
	p.Name = req.Name
	p.Price = req.Price

	if p.Name == "" || p.Price < 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "bad product"})
		return
	}

	p = storeProductAdd(p)
	writeJSON(w, http.StatusCreated, p)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	_ = enc.Encode(v)
}
