package routes

import (
	"html/template"
	"net/http"

	"github.com/alcb1310/bca-go-w-test/database"
	"github.com/alcb1310/bca-go-w-test/utils"
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

	t.budgetRoutes()
	t.HandleFunc("/factura", t.handleInvoices)
	t.HandleFunc("/cierre", t.handleClosure)
}

func (t *transactionRouter) handleInvoices(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)
	type Ret struct {
		UserName string
		Title    string
		Links    utils.LinksType
	}
	retData := Ret{
		UserName: ctxPayload.Name,
		Title:    "BCA - Transacciones",
		Links:    *utils.Links,
	}

	switch r.Method {
	case http.MethodGet:
		file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/transactions/invoice.html")
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

func (t *transactionRouter) handleClosure(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)
	type Ret struct {
		UserName string
		Title    string
		Links    utils.LinksType
	}
	retData := Ret{
		UserName: ctxPayload.Name,
		Title:    "BCA - Transacciones",
		Links:    *utils.Links,
	}

	switch r.Method {
	case http.MethodGet:
		file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/transactions/closure.html")
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
