package routes

import (
	"net/http"
	"text/template"

	"github.com/alcb1310/bca-go-w-test/database"
	"github.com/alcb1310/bca-go-w-test/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *ProtectedRouter) handleUsers(w http.ResponseWriter, r *http.Request) {
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

		http.Redirect(w, r, "/api/v1/edit-user", http.StatusPermanentRedirect)
	case http.MethodGet:
		tmpl, err := template.ParseFiles("templates/bca/users/all-users.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusTeapot)
			return
		}

		users, err := s.db.GetAllUsers(ctxPayload.CompanyId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, users)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *ProtectedRouter) handleSimpleUser(w http.ResponseWriter, r *http.Request) {
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

		sql := "UPDATE \"user\" SET email=$3, name=$4, role_id=$5 WHERE id=$1 AND company_id = $2"
		if _, err := s.db.Exec(sql, userId, ctxPayload.CompanyId, u.Email, u.Name, u.Role); err != nil {
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
		var oldPassword *string
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

		sql := "SELECT password FROM \"user\" WHERE id=$1 and company_id = $2"
		rows, err := s.db.Query(sql, userId, ctxPayload.CompanyId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		defer rows.Close()
		var prev *string
		prev = nil

		for rows.Next() {
			var passwd string
			if err := rows.Scan(&passwd); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			prev = &passwd
		}
		if _, err := utils.ComparePassword(*prev, *oldPassword); err != nil {
			http.Error(w, "Contraseña inválida", http.StatusBadRequest)
			return
		}

		var newPassword *string
		if r.FormValue("new") == "" {
			newPassword = nil
		} else {
			pass := r.FormValue("new")
			newPassword = &pass
		}

		savePasswd, err := utils.EncryptPasssword(*newPassword)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sql = "UPDATE \"user\" SET password = $3 WHERE id=$1 and company_id = $2"
		if _, err := s.db.Exec(sql, userId, ctxPayload.CompanyId, savePasswd); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/api/v1/edit-user", http.StatusPermanentRedirect)
	case http.MethodDelete:
		if ctxPayload.Id == userId {
			http.Error(w, "Can't delete yourself", http.StatusBadRequest)
			return
		}

		if err := s.db.DeleteUser(userId, ctxPayload.CompanyId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/api/v1/edit-user", http.StatusPermanentRedirect)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func (s *ProtectedRouter) tmplAddUser(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/bca/users/add-user.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusTeapot)
		return
	}

	tmpl.Execute(w, nil)
}

func (s *ProtectedRouter) tmplChangePassword(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)
	tmpl, err := template.ParseFiles("templates/bca/users/change-password.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusTeapot)
		return
	}

	tmpl.Execute(w, ctxPayload)
}
