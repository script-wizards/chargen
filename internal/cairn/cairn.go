package cairn

import (
	"net/http"
	"text/template"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/cairn.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
