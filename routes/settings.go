package routes

import (
	"html/template"
	"net/http"

	"github.com/alcb1310/bca-go-w-test/database"
	"github.com/alcb1310/bca-go-w-test/utils"
	"github.com/gorilla/mux"
)

type settingsRouter struct {
	*mux.Router
	db *database.Database
}

func (p *ProtectedRouter) settingsRoutes() {
	s := &settingsRouter{
		Router: p.PathPrefix("/parametros").Subrouter(),
		db:     p.db,
	}

	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	s.HandleFunc("/partidas", s.handleBudgetItems)

	s.projectsRoutes()
	s.supplierRoutes()
}

func (s *settingsRouter) handleBudgetItems(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)
	type Ret struct {
		UserName string
		Title    string
		Links    utils.LinksType
	}
	retData := Ret{
		UserName: ctxPayload.Name,
		Title:    "BCA - Parámetros",
		Links:    *utils.Links,
	}

	switch r.Method {
	case http.MethodGet:
		file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/settings/budget-items.html")
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
