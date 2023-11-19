package routes

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/alcb1310/bca-go-w-test/database"
	"github.com/alcb1310/bca-go-w-test/types"
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
	b.HandleFunc("/crear", b.handleCreateBudget)
}

func (b *budgetRouter) handleCreateBudget(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)
	retData := utils.InitializeMap()
	retData["UserName"] = ctxPayload.Name
	retData["Title"] = "BCA - Presupuesto"

	switch r.Method {
	case http.MethodGet:
		file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/transactions/budget/create-budget.html")
		tmpl, err := template.ParseFiles(file...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusTeapot)
			return
		}

		tmpl.ExecuteTemplate(w, "base", retData)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (b *budgetRouter) handleBudget(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)
	retData := utils.InitializeMap()
	retData["UserName"] = ctxPayload.Name
	retData["Title"] = "BCA - Transacciones"

	switch r.Method {
	case http.MethodGet:
		file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/transactions/budget/index.html")
		tmpl, err := template.ParseFiles(file...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusTeapot)
			return
		}

		pag := &types.PaginationQuery{}
		strPage := r.URL.Query().Get("pagina")
		if strPage == "" {
			pag.Offset = uint(1)
		} else {
			page, err := strconv.Atoi(strPage)
			if err != nil {
				retData["Error"] = err.Error()
			} else {
				pag.Offset = uint(page)
			}
		}
		strLimit := r.URL.Query().Get("items")
		if strLimit == "" {
			pag.Limit = uint(10)
		} else {
			limit, err := strconv.Atoi(strLimit)
			if err != nil {
				retData["Error"] = err.Error()
			} else {
				pag.Limit = uint(limit)
			}
		}

		retData["Budgets"], _, err = b.db.GetBudgets(ctxPayload.CompanyId, pag, "")
		if err != nil {
			retData["Error"] = err.Error()
		}

		tmpl.ExecuteTemplate(w, "base", retData)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
