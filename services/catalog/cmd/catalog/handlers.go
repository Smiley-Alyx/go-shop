package main

import (
	"encoding/json"
	"net/http"
)

func (a app) handleProducts(w http.ResponseWriter, r *http.Request) {
	list := a.store.List()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(list)
}
