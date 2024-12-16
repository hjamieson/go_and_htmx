package main

import (
	"contacts2/store"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
)

const port = ":3000"

var contacts store.Store = store.NewContacts()

type Person = store.Contact

var templates = template.Must(template.ParseGlob("static/*.html"))

func main() {

	mockData()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/contacts", http.StatusFound)
	})
	http.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/contacts", http.StatusFound)
	})
	http.HandleFunc("GET /contacts/{id}", viewContact)
	http.HandleFunc(("GET /contacts/new"), handleNew)
	http.HandleFunc(("POST /contacts/new"), handleSave)
	http.HandleFunc(("GET /contacts/{id}/edit"), handleEdit)
	http.HandleFunc(("POST /contacts/{id}/edit"), handleUpdate)
	http.HandleFunc(("POST /contacts/{id}/delete"), handleDelete)
	http.HandleFunc("GET /contacts", home)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	slog.Info("contacts2", "port", port)
	http.ListenAndServe(port, nil)
}

func mockData() {
	contacts.Add(Person{First: "Hugh", Last: "Janus", Phone: "555-555-5555", Email: "dummy.com"})
	contacts.Add(Person{First: "Mike", Last: "Litoris", Phone: "555-555-5555", Email: "literal@gmail.com"})
	contacts.Add(Person{First: "Al", Last: "Coholic", Phone: "555-555-5555", Email: "beer@gmail.com"})
	contacts.Add(Person{First: "Mike", Last: "Oxmaul", Phone: "555-555-5555", Email: "munt@gmail.com"})
}

func home(w http.ResponseWriter, r *http.Request) {
	slog.Info("home", "method", r.Method)
	params := r.URL.Query()
	slog.Info("query: ", "params", params)
	var result []Person
	if params["q"] != nil {
		result = contacts.Search(params["q"][0])
	} else {
		result = contacts.All()
	}
	slog.Info("home", "defined templates", templates.DefinedTemplates())
	templates.ExecuteTemplate(w, "layout.html", result)
}

func handleNew(w http.ResponseWriter, r *http.Request) {
	slog.Info("new", "url", r.URL.Path, "method", r.Method)
	newGuy := FormContact{Person{}, "", ""}
	templates.ExecuteTemplate(w, "new.html", newGuy)
}

func handleSave(w http.ResponseWriter, r *http.Request) {
	slog.Info("save", "url", r.URL.Path, "method", r.Method)

	newGuy := FormContact{
		Person{First: r.FormValue("first_name"),
			Last:  r.FormValue("last_name"),
			Phone: r.FormValue("phone"),
			Email: r.FormValue("email"),
		}, "", ""}
	if !ValidateContact(&newGuy) {
		slog.Error("new contact fails validation")
		templates.ExecuteTemplate(w, "new.html", newGuy)
		return
	}
	id, _ := contacts.Add(newGuy.Entry)
	slog.Info("contacts", "new", newGuy, "id", id)
	http.Redirect(w, r, "/contacts", http.StatusFound)
}

type FormContact struct {
	Entry      Person
	EmailError string
	PhoneError string
}

func ValidateContact(fc *FormContact) bool {
	if !strings.Contains(fc.Entry.Email, "@") {
		fc.EmailError = "Email must contain @ symbol"
	}
	if len(fc.Entry.Phone) != 12 {
		fc.PhoneError = "Phone number must be in the format 555-555-5555"
	}
	return fc.EmailError == "" && fc.PhoneError == ""
}

func viewContact(w http.ResponseWriter, r *http.Request) {
	slog.Info("view", "url", r.URL.Path, "method", r.Method, "id", r.PathValue("id"))
	id := r.PathValue("id")
	result := contacts.Get(id)
	templates.ExecuteTemplate(w, "view.html", result)
}

func handleEdit(w http.ResponseWriter, r *http.Request) {
	slog.Info("edit", "url", r.URL.Path, "method", r.Method, "id", r.PathValue("id"))
	id := r.PathValue("id")
	var fc FormContact = FormContact{contacts.Get(id), "", ""}
	templates.ExecuteTemplate(w, "edit.html", fc)
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	slog.Info("update", "url", r.URL.Path, "method", r.Method, "id", r.PathValue("id"))
	id := r.PathValue("id")
	current := contacts.Get(id)
	fc := FormContact{
		Person{Id: current.Id,
			First: r.FormValue("first_name"),
			Last:  r.FormValue("last_name"),
			Phone: r.FormValue("phone"),
			Email: r.FormValue("email"),
		}, "", ""}
	if !ValidateContact(&fc) {
		slog.Error("new contact fails validation")
		templates.ExecuteTemplate(w, "edit.html", fc)
		return
	}
	var _ = contacts.Update(id, fc.Entry)
	slog.Info("contacts", "new", fc.Entry, "id", id)
	http.Redirect(w, r, "/contacts", http.StatusFound)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	slog.Info("delete", "url", r.URL.Path, "method", r.Method, "id", r.PathValue("id"))
	id := r.PathValue("id")
	err := contacts.Delete(id)
	if err != nil {
		slog.Error("delete failed", "error", err)
		http.Error(w, "delete failed", http.StatusBadRequest)
		return
	}
	slog.Info("contact deleted", "id", id)
	http.Redirect(w, r, "/contacts", http.StatusFound)
}
