package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	a := app{version: "test"}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", a.handleHealthz)

	r := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d", w.Code)
	}

	if w.Body.String() != "ok\n" {
		t.Fatalf("body=%q", w.Body.String())
	}
}

func TestVersion(t *testing.T) {
	a := app{version: "testver"}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /version", a.handleVersion)

	r := httptest.NewRequest("GET", "/version", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d", w.Code)
	}

	if w.Body.String() != "testver\n" {
		t.Fatalf("body=%q", w.Body.String())
	}
}
