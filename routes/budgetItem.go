package routes

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

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
	b.HandleFunc("/{budgetItemId}", b.handleSingleBudgetItem)
}

func (b *budgetItemRouter) handleSingleBudgetItem(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)
	retData := utils.InitializeMap()
	retData["UserName"] = ctxPayload.Name
	retData["Title"] = "BCA - Parámetros"
	vars := mux.Vars(r)
	budgetItemId, err := uuid.Parse(vars["budgetItemId"])
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	switch r.Method {
	case http.MethodPost:
		bi := &types.BudgetItemType{}
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		code := r.FormValue("code")
		if code == "" {
			bi.Code = nil
		} else {
			bi.Code = &code
		}
		name := r.FormValue("name")
		if name == "" {
			bi.Name = nil
		} else {
			bi.Name = &name
		}
		level := r.FormValue("level")
		if level == "" {
			bi.Level = nil
		} else {
			l, err := strconv.ParseInt(level, 10, 32)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			lNum := uint(l)
			bi.Level = &lNum
		}
		parentId := r.FormValue("parent")
		if parentId == "" {
			bi.ParentId = nil
		} else {
			pId, err := uuid.Parse(parentId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			bi.ParentId = &pId
		}
		referer := r.FormValue("referer")
		accumulates := r.FormValue("accumulate")
		bi.Accumulates = accumulates == "on"
		bi.ID = budgetItemId

		if err := b.db.UpdateBudgetItem(ctxPayload.CompanyId, bi); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		r.Method = http.MethodGet
		http.Redirect(w, r, referer, http.StatusSeeOther)
	case http.MethodGet:
		retData["Referer"] = r.Header.Get("Referer")
		funcs := map[string]any{
			"Compare": strings.Compare,
		}
		file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/settings/budget-items/create-budget-items.html")
		tmpl, err := template.New("base").Funcs(funcs).ParseFiles(file...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusTeapot)
			return
		}
		bi, err := b.db.GetBudgetItemById(ctxPayload.CompanyId, budgetItemId)
		if err != nil {
			retData["Error"] = err
		}
		retData["BudgetItem"] = bi
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
		referer := r.FormValue("referer")

		r.Method = http.MethodGet
		http.Redirect(w, r, referer, http.StatusSeeOther)
	case http.MethodGet:
		file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/settings/budget-items/index.html")
		file = append(file, utils.PaginationTemplate)
		tmpl, err := template.ParseFiles(file...)
		if err != nil {
			retData["Error"] = err.Error()
		}
		retData["Referer"] = r.Header.Get("Referer")

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
		funcs := map[string]any{
			"Compare": strings.Compare,
		}
		tmpl, err := template.New("base").Funcs(funcs).ParseFiles(file...)
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
