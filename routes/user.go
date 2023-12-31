package routes

import (
	"net/http"
	"text/template"

	"github.com/alcb1310/bca-go-w-test/database"
	"github.com/alcb1310/bca-go-w-test/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type usersRouter struct {
	*mux.Router

	db database.Database
}

func (p *ProtectedRouter) usersRoutes() {
	u := &usersRouter{
		Router: p.PathPrefix("/usuarios").Subrouter(),
		db:     *p.db,
	}

	u.HandleFunc("/", u.handleUsers)
	u.HandleFunc("/agregar", u.addUser)
	u.HandleFunc("/contrasena", u.passwordRoute)
	u.HandleFunc("/agregar/{userId}", u.updateUser)
	u.HandleFunc("/{userId}", u.handleSimpleUser)
}

func (s *usersRouter) passwordRoute(w http.ResponseWriter, r *http.Request) {
	file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/users/password.html")
	tmpl, err := template.ParseFiles(file...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusTeapot)
		return
	}

	ctxPayload, _ := getMyPaload(r)
	switch r.Method {
	case http.MethodGet:
		retData := utils.InitializeMap()
		tmpl.ExecuteTemplate(w, "base", retData)

	case http.MethodPost:
		var oldPassword, newPassword *string
		retData := utils.InitializeMap()
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if r.FormValue("old") == "" {
			oldPassword = nil
		} else {
			pass := r.FormValue("old")
			oldPassword = &pass
		}

		if r.FormValue("new") == "" {
			newPassword = nil
		} else {
			pass := r.FormValue("new")
			newPassword = &pass
		}

		if err := s.db.ChangePassword(*oldPassword, *newPassword, ctxPayload.Id, ctxPayload.CompanyId); err != nil {
			retData["Error"] = "Contraseña inválida"
			tmpl.ExecuteTemplate(w, "base", retData)
			return
		}
		r.Method = http.MethodGet
		http.Redirect(w, r, "/bca/usuarios/contrasena", http.StatusSeeOther)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *usersRouter) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := uuid.Parse(vars["userId"])
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	ctxPayload, _ := getMyPaload(r)

	switch r.Method {
	case http.MethodPost:
		u := &database.UserInfo{}
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if r.FormValue("email") == "" {
			u.Email = nil
		} else {
			email := r.FormValue("email")
			u.Email = &email
		}

		if r.FormValue("name") == "" {
			u.Name = nil
		} else {
			name := r.FormValue("name")
			u.Name = &name
		}

		if r.FormValue("role") == "" {
			u.Role = nil
		} else {
			role := r.FormValue("role")
			u.Role = &role
		}

		if err := s.db.UpdateUser(u, userId, ctxPayload.CompanyId); err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		r.Method = http.MethodGet
		http.Redirect(w, r, "/bca/usuarios/", http.StatusSeeOther)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *usersRouter) addUser(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)
	retData := utils.InitializeMap()

	file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/users/add-user.html")
	tmpl, err := template.ParseFiles(file...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusTeapot)
		return
	}

	retData["UserName"] = ctxPayload.Name
	retData["Title"] = "BCA - Transacciones"

	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		tmpl.ExecuteTemplate(w, "base", retData)
	case http.MethodPost:
		u := &database.UserInfo{}
		// Data sanitization
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if r.FormValue("email") == "" {
			u.Email = nil
		} else {
			email := r.FormValue("email")
			u.Email = &email
		}

		if r.FormValue("name") == "" {
			u.Name = nil
		} else {
			name := r.FormValue("name")
			u.Name = &name
		}

		if r.FormValue("password") == "" {
			u.Password = nil
		} else {
			password := r.FormValue("password")
			u.Password = &password
		}
		if r.FormValue("role") == "" {
			u.Role = nil
		} else {
			role := r.FormValue("role")
			u.Role = &role
		}
		// end of data sanitization

		u.CompanyId = ctxPayload.CompanyId
		if err := s.db.AddUser(*u); err != nil {
			retData["Error"] = err.Error()
			w.WriteHeader(http.StatusConflict)
			tmpl.ExecuteTemplate(w, "base", retData)
			return
		}

		r.Method = http.MethodGet
		http.Redirect(w, r, "/bca/usuarios/", http.StatusSeeOther)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *usersRouter) handleUsers(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)

	retData := utils.InitializeMap()
	retData["UserName"] = ctxPayload.Name
	retData["Title"] = "BCA - Transacciones"
	switch r.Method {
	case http.MethodGet:
		searchParam := r.URL.Query().Get("usuario")
		file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/users/all-users.html")
		tmpl, err := template.ParseFiles(file...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusTeapot)
			return
		}

		retData["Users"], err = s.db.GetAllUsers(ctxPayload.CompanyId, searchParam)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		tmpl.ExecuteTemplate(w, "base", retData)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *usersRouter) handleSimpleUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := uuid.Parse(vars["userId"])
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	ctxPayload, _ := getMyPaload(r)
	retData := utils.InitializeMap()

	switch r.Method {
	case http.MethodGet:
		retData["UserName"] = ctxPayload.Name
		retData["Title"] = "BCA - Transacciones"
		file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/users/add-user.html")
		tmpl, err := template.ParseFiles(file...)
		user, err := s.db.GetOneUser(userId, ctxPayload.CompanyId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusTeapot)
			return
		}
		retData["User"] = user

		tmpl.ExecuteTemplate(w, "base", retData)
	case http.MethodDelete:
		if ctxPayload.Id == userId {
			http.Error(w, "No se puede eliminar a si mismo", http.StatusBadRequest)
			return
		}

		if err := s.db.DeleteUser(userId, ctxPayload.CompanyId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		r.Method = http.MethodGet
		http.Redirect(w, r, "/bca/usuarios/", http.StatusSeeOther)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}
