package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultPort = "8082"
)

type app struct {
	version string
}

func main() {
	a := app{
		version: getenv("VERSION", "0.0.0"),
	}

	storeOrdersInit()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", a.handleHealthz)
	mux.HandleFunc("GET /version", a.handleVersion)
	mux.HandleFunc("POST /orders", a.handleOrdersCreate)
	mux.HandleFunc("GET /orders/{id}", a.handleOrdersGet)
	mux.HandleFunc("GET /orders/{id}/status", a.handleOrdersGetStatus)
	mux.HandleFunc("POST /orders/{id}/status", a.handleOrdersSetStatus)

	srv := &http.Server{
		Addr:              ":" + getenv("PORT", defaultPort),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("order: listen %s version=%s", srv.Addr, a.version)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("order: listen: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("order: shutdown")
	err := srv.Shutdown(shutdownCtx)
	if err != nil {
		log.Printf("order: shutdown error: %v", err)
	}
}

func (a app) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
}

func (a app) handleVersion(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "%s\n", a.version)
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
