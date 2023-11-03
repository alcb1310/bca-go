package routes

import (
	"html/template"
	"net/http"

	"github.com/alcb1310/bca-go-w-test/database"
	"github.com/alcb1310/bca-go-w-test/types"
	"github.com/alcb1310/bca-go-w-test/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type proyectRouter struct {
	*mux.Router

	db *database.Database
}

func (s *settingsRouter) projectsRoutes() {
	p := &proyectRouter{
		Router: s.PathPrefix("/proyectos").Subrouter(),
		db:     s.db,
	}

	p.HandleFunc("/", p.handleProjects)
	p.HandleFunc("/crear", p.handleCreateProject)
	p.HandleFunc("/{userId}", p.handleEditProject)
}

func (p *proyectRouter) handleEditProject(w http.ResponseWriter, r *http.Request) {
	file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/settings/projects/create-project.html")
	tmpl, err := template.ParseFiles(file...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusTeapot)
		return
	}

	vars := mux.Vars(r)
	userId, err := uuid.Parse(vars["userId"])
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	ctxPayload, _ := getMyPaload(r)
	retData := utils.InitializeMap()
	retData["UserName"] = ctxPayload.Name
	retData["Title"] = "BCA - Transacciones"

	switch r.Method {
	case http.MethodGet:
		project, err := p.db.GetSingleProject(userId, ctxPayload.CompanyId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusTeapot)
			return
		}
		if project.Name == nil {
			http.Error(w, "Proyecto inexistente", http.StatusNotFound)
			return
		}

		retData["Project"] = project
		tmpl.ExecuteTemplate(w, "base", retData)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (p *proyectRouter) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)
	retData := utils.InitializeMap()
	retData["UserName"] = ctxPayload.Name
	retData["Title"] = "BCA - Transacciones"

	switch r.Method {
	case http.MethodGet:
		file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/settings/projects/create-project.html")
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

func (p *proyectRouter) handleProjects(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)
	retData := utils.InitializeMap()
	retData["UserName"] = ctxPayload.Name
	retData["Title"] = "BCA - Transacciones"

	switch r.Method {
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		project := &types.Project{
			CompanyId: ctxPayload.CompanyId,
		}
		name := r.FormValue("name")
		active := r.PostFormValue("active")
		if name == "" {
			project.Name = nil
		} else {
			project.Name = &name
		}

		if active == "on" {
			project.IsActive = true
		} else {
			project.IsActive = false
		}

		if err := p.db.CreateProject(project); err != nil {
			retData["Error"] = err.Error()
			file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/settings/projects/create-project.html")
			tmpl, err := template.ParseFiles(file...)
			if err != nil {
				http.Error(w, err.Error(), http.StatusTeapot)
				return
			}

			tmpl.ExecuteTemplate(w, "base", retData)
			return
		}

		r.Method = http.MethodGet
		http.Redirect(w, r, "/bca/parametros/proyectos/", http.StatusSeeOther)
	case http.MethodGet:
		file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/settings/projects.html")
		tmpl, err := template.ParseFiles(file...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusTeapot)
			return
		}

		retData["Projects"], err = p.db.GetAllProjects(ctxPayload.CompanyId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "base", retData)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
