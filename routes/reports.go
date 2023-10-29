package routes

import (
	"net/http"

	"github.com/alcb1310/bca-go-w-test/database"
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
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Actual page"))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (t *reportsRouter) handleHistoric(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Historico page"))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
