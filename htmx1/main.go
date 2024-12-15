package main

import (
	"html/template"
	"log"
	"log/slog"
	"net/http"
)

func main() {
	slog.Info("Starting server on http://localhost:3000")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("."))))

	templates := template.Must(template.ParseGlob("*.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/form1", http.StatusSeeOther)
	})

	http.Handle("/contacts", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.FormValue("q")
		code := r.FormValue("code")
		coupon := r.FormValue("coupon")
		slog.Info("calling contacts", "q", query, "code", code, "coupon", coupon)
		http.ServeFile(w, r, "contacts.html")
	}))

	http.HandleFunc("GET /form1", func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "form1.html", nil)
	})
	http.HandleFunc("/form2", func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "form2.html", nil)
	})
	http.HandleFunc("/form3", func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "form3.html", nil)
	})

	http.HandleFunc("POST /form1", form1Post)

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func form1Post(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.FormValue("q")
	slog.Info("calling form1", "method", r.Method)
	txt := `<ol><li>{{.}}</li><li>{{.}}</li><li>{{.}}</li></ol>`
	tmpl, err := template.New("list").Parse(txt)
	if err != nil {
		slog.Error("template error:", "message", err.Error())
		w.Write([]byte("<p>Error: " + err.Error() + "</p>"))
		return
	}
	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, searchTerm)
}
