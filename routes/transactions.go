package routes

import (
	"net/http"

	"github.com/alcb1310/bca-go-w-test/database"
	"github.com/gorilla/mux"
)

type transactionRouter struct {
	*mux.Router

	db *database.Database
}

func (p *ProtectedRouter) transactionsRoutes() {
	t := &transactionRouter{
		Router: p.PathPrefix("/transacciones").Subrouter(),
		db:     p.db,
	}

	t.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	t.HandleFunc("/presupuesto", t.handleBudget)
	t.HandleFunc("/factura", t.handleInvoices)
	t.HandleFunc("/cierre", t.handleClosure)
}

func (t *transactionRouter) handleBudget(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Presupuesto page"))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (t *transactionRouter) handleInvoices(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Facturas page"))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (t *transactionRouter) handleClosure(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Cierre page"))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
