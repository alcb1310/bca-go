package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/alcb1310/bca-go-w-test/database"
	"github.com/alcb1310/bca-go-w-test/types"
	"github.com/alcb1310/bca-go-w-test/utils"
	"github.com/google/uuid"
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
	s.HandleFunc("/agregar", s.createSupplier)
	s.HandleFunc("/{supplierId}", s.handelSingleSupplier)
}

func (s *supplierRouter) handelSingleSupplier(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	supplierId, err := uuid.Parse(vars["supplierId"])
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	ctxPayload, _ := getMyPaload(r)
	retData := utils.InitializeMap()
	retData["UserName"] = ctxPayload.Name
	retData["Title"] = "BCA - Parámetros"
	file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/settings/suppliers/create-supplier.html")
	tmpl, err := template.ParseFiles(file...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusTeapot)
		return
	}

	switch r.Method {
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			retData["Error"] = err.Error()
			tmpl.ExecuteTemplate(w, "base", retData)
			return
		}
		sup := &types.SupplierType{
			CompanyId: ctxPayload.CompanyId,
			ID:        supplierId,
		}

		ruc := r.PostFormValue("ruc")
		name := r.FormValue("name")
		contactName := r.FormValue("contact_name")
		contactEmail := r.FormValue("contact_email")
		contactPhone := r.FormValue("contact_phone")

		if ruc == "" {
			sup.Ruc = nil
		} else {
			sup.Ruc = &ruc
		}
		if name == "" {
			sup.Name = nil
		} else {
			sup.Name = &name
		}
		if contactName == "" {
			sup.ContactName = nil
		} else {
			sup.ContactName = &contactName
		}
		if contactEmail == "" {
			sup.ContactEmail = nil
		} else {
			sup.ContactEmail = &contactEmail
		}
		if contactPhone == "" {
			sup.ContactPhone = nil
		} else {
			sup.ContactPhone = &contactPhone
		}

		if err := s.db.UpdateSupplier(sup); err != nil {
			retData["Error"] = err.Error()
			tmpl.ExecuteTemplate(w, "base", retData)
			return
		}

		r.Method = http.MethodGet
		http.Redirect(w, r, "/bca/parametros/proveedor/", http.StatusSeeOther)
	case http.MethodGet:
		sup, err := s.db.GetSingleSupplier(supplierId, ctxPayload.CompanyId)
		if err != nil {
			retData["Error"] = err.Error()
		}
		retData["Supplier"] = sup
		tmpl.ExecuteTemplate(w, "base", retData)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *supplierRouter) createSupplier(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)
	retData := utils.InitializeMap()
	retData["UserName"] = ctxPayload.Name
	retData["Title"] = "BCA - Parámetros"

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
	retData := utils.InitializeMap()
	retData["UserName"] = ctxPayload.Name
	retData["Title"] = "BCA - Parámetros"

	switch r.Method {
	case http.MethodPost:
		file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/settings/suppliers/create-supplier.html")
		tmpl, err := template.ParseFiles(file...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusTeapot)
			return
		}

		if err := r.ParseForm(); err != nil {
			retData["Error"] = err.Error()
			tmpl.ExecuteTemplate(w, "base", retData)
			return
		}
		sup := &types.SupplierType{
			CompanyId: ctxPayload.CompanyId,
		}

		ruc := r.PostFormValue("ruc")
		name := r.FormValue("name")
		contactName := r.FormValue("contact_name")
		contactEmail := r.FormValue("contact_email")
		contactPhone := r.FormValue("contact_phone")

		if ruc == "" {
			sup.Ruc = nil
		} else {
			sup.Ruc = &ruc
		}
		if name == "" {
			sup.Name = nil
		} else {
			sup.Name = &name
		}
		if contactName == "" {
			sup.ContactName = nil
		} else {
			sup.ContactName = &contactName
		}
		if contactEmail == "" {
			sup.ContactEmail = nil
		} else {
			sup.ContactEmail = &contactEmail
		}
		if contactPhone == "" {
			sup.ContactPhone = nil
		} else {
			sup.ContactPhone = &contactPhone
		}

		if err := s.db.CreateSupplier(sup); err != nil {
			retData["Error"] = err.Error()
			tmpl.ExecuteTemplate(w, "base", retData)
			return
		}

		r.Method = http.MethodGet
		http.Redirect(w, r, "/bca/parametros/proveedor/", http.StatusSeeOther)
	case http.MethodGet:
		file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/settings/suppliers/suppliers.html")
		file = append(file, utils.PaginationTemplate)
		tmpl, err := template.ParseFiles(file...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(fmt.Sprintf("Error parsing template: %s", err.Error()))
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

		sup, pagin, err := s.db.GetAllSuppliers(ctxPayload.CompanyId, pag)
		if err != nil {
			retData["Error"] = err.Error()
		}
		retData["Suppliers"] = sup
		retData["Pagination"] = pagin
		retData["URL"] = "/bca/parametros/proveedor/"

		tmpl.ExecuteTemplate(w, "base", retData)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
