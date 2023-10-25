package routes

import (
	"fmt"
	"net/http"
	"text/template"
)

func (p *ProtectedRouter) handleProject(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)
	switch r.Method {
	case http.MethodGet:
		tmpl, err := template.ParseFiles("templates/bca/projects/all-projects.html")
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusTeapot)
			return
		}

		tmpl.Execute(w, ctxPayload)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
