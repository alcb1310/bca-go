package routes

import (
	"html/template"
	"net/http"

	"github.com/alcb1310/bca-go-w-test/database"
	"github.com/alcb1310/bca-go-w-test/utils"
	"github.com/gorilla/mux"
)

type supplierRouter struct {
	*mux.Router

	db *database.Database
}

func (r *settingsRouter) supplierRoutes() {
	s := &supplierRouter{
		Router: r.PathPrefix("/proveedor").Subrouter(),
		db:     r.db,
	}

	s.HandleFunc("/", s.handleSuppliers)
	s.HandleFunc("/crear", s.createSupplier)
}

func (s *supplierRouter) createSupplier(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)
	retData := make(map[string]interface{})
	retData["UserName"] = ctxPayload.Name
	retData["Title"] = "BCA - Parámetros"
	retData["Links"] = *utils.Links

	switch r.Method {
	case http.MethodGet:
		file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/settings/suppliers/create-supplier.html")
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

func (s *supplierRouter) handleSuppliers(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)
	retData := make(map[string]interface{})
	retData["UserName"] = ctxPayload.Name
	retData["Title"] = "BCA - Parámetros"
	retData["Links"] = *utils.Links

	switch r.Method {
	case http.MethodPost:
		file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/settings/suppliers/create-supplier.html")
		tmpl, err := template.ParseFiles(file...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusTeapot)
			return
		}

		tmpl.ExecuteTemplate(w, "base", retData)
	case http.MethodGet:
		file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/settings/suppliers/suppliers.html")
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
