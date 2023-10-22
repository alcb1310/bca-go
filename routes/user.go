package routes

import (
	"net/http"
	"strings"
	"text/template"

	"github.com/alcb1310/bca-go-w-test/types"
	"github.com/alcb1310/bca-go-w-test/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UserInfo struct {
	email    *string
	name     *string
	password *string
	role     *string
}

func (s *ProtectedRouter) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		u := &UserInfo{}
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if r.FormValue("email") == "" {
			u.email = nil
		} else {
			email := r.FormValue("email")
			u.email = &email
		}

		if r.FormValue("name") == "" {
			u.name = nil
		} else {
			name := r.FormValue("name")
			u.name = &name
		}

		if r.FormValue("password") == "" {
			u.password = nil
		} else {
			password := r.FormValue("password")
			encryptedPassword, err := utils.EncryptPasssword(password)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			password = string(encryptedPassword)
			u.password = &password
		}

		if r.FormValue("role") == "" {
			u.password = nil
		} else {
			role := r.FormValue("role")
			u.role = &role
		}
		ctxPayload, _ := getMyPaload(r)
		sql := "INSERT INTO \"user\" (email, name, password, company_id, role_id) VALUES($1, $2, $3, $4, $5)"
		if _, err := s.db.Exec(sql, &u.email, &u.name, &u.password, ctxPayload.CompanyId, &u.role); err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		http.Redirect(w, r, "/api/v1/edit-user", http.StatusPermanentRedirect)
	case http.MethodGet:
		// if r.Method == http.MethodGet {
		ctxPayload, _ := getMyPaload(r)
		tmpl, err := template.ParseFiles("templates/bca/users/all-users.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusTeapot)
			return
		}

		sql := "SELECT user_id, user_email, user_name, role_name FROM user_without_password WHERE company_id=$1"
		rows, err := s.db.Query(sql, ctxPayload.CompanyId)
		if err != nil {
			http.Error(w, "Error al buscar usuarios", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var Users []types.User

		for rows.Next() {
			var id, email, name, role string
			if err := rows.Scan(&id, &email, &name, &role); err != nil {
				http.Error(w, "Error al buscar usuarios", http.StatusInternalServerError)
				return
			}

			strUUID, err := uuid.Parse(id)
			if err != nil {
				http.Error(w, "Error al buscar usuarios", http.StatusInternalServerError)
				return
			}

			Users = append(Users, types.User{
				Id:     strUUID,
				Email:  email,
				Name:   name,
				RoleId: role,
			})
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, Users)
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
	case http.MethodGet:
		sql := "SELECT user_id, user_email, user_name, role_id FROM user_without_password where company_id = $2 and user_id = $1"
		// sql := "SELECT user_id, user_email, user_name, role_name FROM user_without_password WHERE company_id=$1"
		rows, err := s.db.Query(sql, userId, ctxPayload.CompanyId)
		if err != nil {
			http.Error(w, "Error al buscar usuarios", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		user := types.User{}

		for rows.Next() {
			var id, email, name, roleId string
			if err := rows.Scan(&id, &email, &name, &roleId); err != nil {
				http.Error(w, "Error al buscar usuarios", http.StatusInternalServerError)
				return
			}
			foundUserId, err := uuid.Parse(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			user.Id = foundUserId
			user.Name = name
			user.Email = email
			user.CompanyId = ctxPayload.CompanyId
			user.RoleId = strings.Trim(roleId, " ")
		}
		tmpl, err := template.ParseFiles("templates/bca/users/add-user.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusTeapot)
			return
		}

		tmpl.Execute(w, user)
	case http.MethodDelete:
		if ctxPayload.Id == userId {
			http.Error(w, "Can't delete yourself", http.StatusBadRequest)
			return
		}

		sql := "DELETE FROM \"user\" WHERE id = $1 and company_id = $2"
		_ = sql
		if _, err := s.db.Exec(sql, userId, ctxPayload.CompanyId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/api/v1/edit-user", http.StatusPermanentRedirect)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func (s *ProtectedRouter) tmplAddUser(w http.ResponseWriter, r *http.Request) {
	_, _ = getMyPaload(r)
	tmpl, err := template.ParseFiles("templates/bca/users/add-user.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusTeapot)
		return
	}

	tmpl.Execute(w, nil)
}