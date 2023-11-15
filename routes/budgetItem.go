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
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		budgetItem := &types.BudgetItemType{}
		code := r.FormValue("code")
		if code == "" {
			budgetItem.Code = nil
		} else {
			budgetItem.Code = &code
		}
		name := r.FormValue("name")
		if name == "" {
			budgetItem.Name = nil
		} else {
			budgetItem.Name = &name
		}
		level := r.FormValue("level")
		if level == "" {
			budgetItem.Level = nil
		} else {
			levelInt, _ := strconv.Atoi(level)
			levelUint := uint(levelInt)
			budgetItem.Level = &levelUint
		}
		if r.FormValue("accumulate") == "on" {
			budgetItem.Accumulates = true
		} else {
			budgetItem.Accumulates = false
		}
		if r.FormValue("parent") == "" {
			budgetItem.ParentId = nil
		} else {
			parentId, _ := uuid.Parse(r.FormValue("parent"))
			budgetItem.ParentId = &parentId
		}

		if err := b.db.CreateBudgetItem(ctxPayload.CompanyId, budgetItem); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		r.Method = http.MethodGet
		http.Redirect(w, r, "/bca/parametros/partidas/", http.StatusSeeOther)
	case http.MethodGet:
		file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/settings/budget-items/index.html")
		file = append(file, utils.PaginationTemplate)
		tmpl, err := template.ParseFiles(file...)
		if err != nil {
			retData["Error"] = err.Error()
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

		searchParam := r.URL.Query().Get("partida")

		budgetItems, pagin, err := b.db.GetAllBudgetItems(ctxPayload.CompanyId, pag, searchParam)
		if err != nil {
			retData["Error"] = err.Error()
		}
		retData["BudgetItems"] = budgetItems
		retData["Pagination"] = pagin
		if searchParam != "" {
			retData["URL"] = "/bca/parametros/partidas/" + "?partida=" + searchParam + "&"
		} else {
			retData["URL"] = "/bca/parametros/partidas/" + "?"
		}

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