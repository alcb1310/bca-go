package routes

import (
	"html/template"
	"net/http"

	"github.com/alcb1310/bca-go-w-test/database"
	"github.com/alcb1310/bca-go-w-test/utils"
	"github.com/gorilla/mux"
)

type reportsRouter struct {
	*mux.Router

	db *database.Database
}

func (p *ProtectedRouter) reportsRoutes() {
	t := &reportsRouter{
		Router: p.PathPrefix("/reportes").Subrouter(),
		db:     p.db,
	}

	t.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	t.HandleFunc("/actual", t.handleActual)
	t.HandleFunc("/historico", t.handleHistoric)
}

func (t *reportsRouter) handleActual(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)
	type Ret struct {
		UserName string
		Title    string
		Links    utils.LinksType
	}
	retData := Ret{
		UserName: ctxPayload.Name,
		Title:    "BCA - Reportes",
		Links:    *utils.Links,
	}

	switch r.Method {
	case http.MethodGet:
		file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/reports/actual.html")
		tmpl, err := template.ParseFiles(file...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusTeapot)
			return
		}

		tmpl.ExecuteTemplate(w, "base", retData)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (t *reportsRouter) handleHistoric(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)
	type Ret struct {
		UserName string
		Title    string
		Links    utils.LinksType
	}
	retData := Ret{
		UserName: ctxPayload.Name,
		Title:    "BCA - Reportes",
		Links:    *utils.Links,
	}

	switch r.Method {
	case http.MethodGet:
		file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/reports/historic.html")
		tmpl, err := template.ParseFiles(file...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusTeapot)
			return
		}

		tmpl.ExecuteTemplate(w, "base", retData)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
