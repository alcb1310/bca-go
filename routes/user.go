package routes

import (
	"net/http"
	"text/template"

	"github.com/alcb1310/bca-go-w-test/types"
	"github.com/google/uuid"
)

func (s *ProtectedRouter) handleUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
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
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}
