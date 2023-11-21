package routes

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/alcb1310/bca-go-w-test/database"
	"github.com/alcb1310/bca-go-w-test/types"
	"github.com/alcb1310/bca-go-w-test/utils"
	"github.com/google/uuid"
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

		retData["Projects"], err = b.db.GetActiveProjects(ctxPayload.CompanyId)
		if err != nil {
			retData["Error"] = err.Error()
		}

		retData["BudgetItemList"], err = b.db.AllBudgetItemsByAccumulates(ctxPayload.CompanyId, false)
		if err != nil {
			retData["Error"] = err.Error()
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
	case http.MethodPost:
		file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/transactions/budget/create-budget.html")
		tmpl, err := template.ParseFiles(file...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusTeapot)
			return
		}
		retData["Projects"], err = b.db.GetActiveProjects(ctxPayload.CompanyId)
		if err != nil {
			retData["Error"] = err.Error()
			tmpl.ExecuteTemplate(w, "base", retData)
			return
		}

		retData["BudgetItemList"], err = b.db.AllBudgetItemsByAccumulates(ctxPayload.CompanyId, false)
		if err != nil {
			retData["Error"] = err.Error()
			tmpl.ExecuteTemplate(w, "base", retData)
			return
		}

		if err := r.ParseForm(); err != nil {
			retData["Error"] = err.Error()
			tmpl.ExecuteTemplate(w, "base", retData)
			return
		}
		budget := &types.BudgetCreate{}

		project := r.FormValue("project")
		if project == "" {
			retData["Error"] = "Ingrese un proyecto"
			tmpl.ExecuteTemplate(w, "base", retData)
			return
		} else {
			projectUUID, err := uuid.Parse(project)
			if err != nil {
				retData["Error"] = err.Error()
				tmpl.ExecuteTemplate(w, "base", retData)
				return
			}
			budget.ProjectId = &projectUUID
		}

		budgetItem := r.FormValue("budgetItem")
		if budgetItem == "" {
			retData["Error"] = "Ingrese una partida"
			tmpl.ExecuteTemplate(w, "base", retData)
			return
		} else {
			budgetItemUUID, err := uuid.Parse(budgetItem)
			if err != nil {
				retData["Error"] = err.Error()
				tmpl.ExecuteTemplate(w, "base", retData)
				return
			}
			budget.BudgetItemId = &budgetItemUUID
		}

		quantity, err := strconv.ParseFloat(r.FormValue("quantity"), 64)
		if err != nil {
			retData["Error"] = err.Error()
			tmpl.ExecuteTemplate(w, "base", retData)
			return
		}
		budget.Quantity = &quantity

		cost, err := strconv.ParseFloat(r.FormValue("cost"), 64)
		if err != nil {
			retData["Error"] = err.Error()
			tmpl.ExecuteTemplate(w, "base", retData)
			return
		}
		budget.Cost = &cost

		total := quantity * cost
		budget.Total = &total

		err = b.db.CreateBudget(budget, ctxPayload.CompanyId)
		if err != nil {
			retData["Error"] = err.Error()
			tmpl.ExecuteTemplate(w, "base", retData)
			return
		}

		r.Method = http.MethodGet
		http.Redirect(w, r, "/bca/transacciones/presupuesto/", http.StatusSeeOther)
	case http.MethodGet:
		funcMap := template.FuncMap{
			"FormatNumber": utils.FormatNumber,
		}
		file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/transactions/budget/index.html")
		tmpl, err := template.New("base").Funcs(funcMap).ParseFiles(file...)
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
