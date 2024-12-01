package main

import (
	"log/slog"
	"net/http"
)

func main() {
	http.HandleFunc("GET /click", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("/click!")
		w.Write([]byte(`<h1>Clicked!</h1>`))
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	slog.Info("Server", "port", 3000)
	http.ListenAndServe(":3000", nil)
}
