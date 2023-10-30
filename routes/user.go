package routes

import (
	"net/http"
	"text/template"

	"github.com/alcb1310/bca-go-w-test/database"
	"github.com/alcb1310/bca-go-w-test/types"
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
	u.HandleFunc("/{userId}", u.handleSimpleUser)
}

func (s *usersRouter) addUser(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)

	switch r.Method {
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
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		http.Redirect(w, r, "/bca/usuarios/", http.StatusPermanentRedirect)
	case http.MethodGet:
		file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/users/add-user.html")
		tmpl, err := template.ParseFiles(file...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusTeapot)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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

		w.WriteHeader(http.StatusOK)
		tmpl.ExecuteTemplate(w, "base", retData)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *usersRouter) handleUsers(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)
	file := append(utils.RequiredFiles, utils.TEMPLATE_DIR+"bca/users/all-users.html")
	tmpl, err := template.ParseFiles(file...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusTeapot)
		return
	}

	users, err := s.db.GetAllUsers(ctxPayload.CompanyId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type Ret struct {
		UserName string
		Title    string
		Links    utils.LinksType
		Users    []types.User
	}
	retData := Ret{
		UserName: ctxPayload.Name,
		Title:    "BCA - Transacciones",
		Links:    *utils.Links,
		Users:    users,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.ExecuteTemplate(w, "base", retData)
	// }
}

func (s *usersRouter) handleSimpleUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := uuid.Parse(vars["userId"])
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	ctxPayload, _ := getMyPaload(r)

	switch r.Method {
	case http.MethodPut:
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

		http.Redirect(w, r, "/api/v1/edit-user", http.StatusPermanentRedirect)
	case http.MethodGet:
		user, err := s.db.GetOneUser(userId, ctxPayload.CompanyId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusTeapot)
			return
		}

		tmpl, err := template.ParseFiles("templates/bca/users/add-user.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusTeapot)
			return
		}

		tmpl.Execute(w, user)
	case http.MethodPatch:
		var oldPassword, newPassword *string
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

		if err := s.db.ChangePassword(*oldPassword, *newPassword, userId, ctxPayload.CompanyId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/bca/usuarios/", http.StatusPermanentRedirect)
	case http.MethodDelete:
		if ctxPayload.Id == userId {
			http.Error(w, "No se puede eliminar a si mismo", http.StatusBadRequest)
			return
		}

		if err := s.db.DeleteUser(userId, ctxPayload.CompanyId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/bca/usuarios", http.StatusPermanentRedirect)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func (s *usersRouter) tmplChangePassword(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)
	tmpl, err := template.ParseFiles("templates/bca/users/change-password.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusTeapot)
		return
	}

	tmpl.Execute(w, ctxPayload)
}
