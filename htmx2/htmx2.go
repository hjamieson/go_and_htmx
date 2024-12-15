package main

import (
	"html/template"
	"log"
	"log/slog"
	"net/http"
)

const serverPort = ":3000"

var templates *template.Template

func main() {
	templates = template.Must(template.ParseGlob("template/*"))
	log.Println("templates", templates.DefinedTemplates())
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/hello", hello)
	mux.HandleFunc("/index", handleIndex)
	mux.HandleFunc("/boost", handleBoost)
	mux.HandleFunc("/boosted", handleBoosted)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/boostedForm", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/boostedForm.html")
	})
	mux.HandleFunc("/formBoost", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("formBoost", "password", r.FormValue("password"))
		w.Write([]byte("<p>response from server</p>"))
	})
	
	mux.HandleFunc("/prefix/{arg}", handlePrefix)

	s := &http.Server{
		Addr:    serverPort,
		Handler: mux,
	}
	slog.Info("server started", "port", serverPort)
	log.Fatal(s.ListenAndServe())
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	slog.Info("index", "request", r.URL.Path)
	templates.ExecuteTemplate(w, "index.html", nil)
}

func handleBoost(w http.ResponseWriter, r *http.Request) {
	slog.Info("boost", "request", r.URL.Path)
	templates.ExecuteTemplate(w, "boost.html", nil)
}

func handleBoosted(w http.ResponseWriter, r *http.Request) {
	slog.Info("boosted", "request", r.URL.Path)
	text := `<div id="boosted" hx-get="/boosted" hx-trigger="click">Click me</div>`
	w.Write([]byte(text))
}

func handlePrefix(w http.ResponseWriter, r *http.Request) {
	arg := r.PathValue("arg")
	result := `<p class="arg">` + arg + `</p>`
	w.Write([]byte(result))
}
