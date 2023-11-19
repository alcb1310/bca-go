package routes

import (
	"html/template"
	"net/http"

	"github.com/alcb1310/bca-go-w-test/database"
	"github.com/alcb1310/bca-go-w-test/utils"
	"github.com/gorilla/mux"
)

type budgetRouter struct {
	*mux.Router

	db *database.Database
}

func (t *transactionRouter) budgetRoutes() {
	b := &budgetRouter{
		Router: t.PathPrefix("/presupuesto").Subrouter(),
		db:     t.db,
	}

	b.HandleFunc("/", b.handleBudget)
}

func (b *budgetRouter) handleBudget(w http.ResponseWriter, r *http.Request) {
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
		file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/transactions/budget.html")
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
