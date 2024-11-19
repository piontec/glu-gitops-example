package main

import (
	"log/slog"
	"net/http"
)

func main() {
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html><body><h1>Welcome!</h1></body></html>`))
	}))

	slog.Info("Starting server", "addr", ":8080")

	http.ListenAndServe(":8080", nil)
}
