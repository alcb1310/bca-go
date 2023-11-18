package routes

import (
	"net/http"

	"github.com/alcb1310/bca-go-w-test/database"
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

	s.projectsRoutes()
	s.supplierRoutes()
	s.budgetItemsRoutes()
}
