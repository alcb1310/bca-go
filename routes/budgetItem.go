package routes

import (
	"html/template"
	"net/http"

	"github.com/alcb1310/bca-go-w-test/database"
	"github.com/alcb1310/bca-go-w-test/utils"
	"github.com/gorilla/mux"
)

type budgetItemRouter struct {
	*mux.Router
	db *database.Database
}

func (s *settingsRouter) budgetItemsRoutes() {
	b := &budgetItemRouter{
		Router: s.PathPrefix("/partidas").Subrouter(),
		db:     s.db,
	}

	b.HandleFunc("/", b.handleBudgetItems)
	b.HandleFunc("/agregar", b.handleCreateBudgetItem)
}

func (b *budgetItemRouter) handleBudgetItems(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)
	retData := utils.InitializeMap()
	retData["UserName"] = ctxPayload.Name
	retData["Title"] = "BCA - Parámetros"

	switch r.Method {
	case http.MethodGet:
		file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/settings/budget-items/index.html")
		tmpl, err := template.ParseFiles(file...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusTeapot)
			return
		}

		budgetItems, err := b.db.GetAllBudgetItems(ctxPayload.CompanyId)
		if err != nil {
			return
		}
		retData["BudgetItems"] = budgetItems

		tmpl.ExecuteTemplate(w, "base", retData)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (b *budgetItemRouter) handleCreateBudgetItem(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)
	retData := utils.InitializeMap()
	retData["UserName"] = ctxPayload.Name
	retData["Title"] = "BCA - Parámetros"
	referer := r.Header.Get("Referer")
	retData["Referer"] = referer
	switch r.Method {
	case http.MethodGet:
		file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/settings/budget-items/create-budget-items.html")
		tmpl, err := template.ParseFiles(file...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusTeapot)
			return
		}

		retData["BudgetItemList"], err = b.db.AllBudgetItemsByAccumulates(ctxPayload.CompanyId, true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusTeapot)
			return
		}

		tmpl.ExecuteTemplate(w, "base", retData)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
