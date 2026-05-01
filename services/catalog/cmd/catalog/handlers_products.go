package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type createProductRequest struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type apiError struct {
	Error   string         `json:"error"`
	Code    string         `json:"code"`
	Details map[string]any `json:"details,omitempty"`
}

func (a app) handleProductsList(w http.ResponseWriter, r *http.Request) {
	items := storeProductList()
	writeJSON(w, http.StatusOK, items)
}

func (a app) handleProductsGet(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad id", Code: "bad_request"})
		return
	}

	p, ok := storeProductGetByID(id)
	if ok == 0 {
		writeJSON(w, http.StatusNotFound, apiError{Error: "not found", Code: "not_found"})
		return
	}

	writeJSON(w, http.StatusOK, p)
}

func (a app) handleProductsCreate(w http.ResponseWriter, r *http.Request) {
	var req createProductRequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad json", Code: "bad_request"})
		return
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad json", Code: "bad_request"})
		return
	}

	p := Product{}
	p.Name = req.Name
	p.Price = req.Price

	if p.Name == "" || p.Price < 0 {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad product", Code: "validation"})
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
