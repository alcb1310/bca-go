package routes

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/alcb1310/bca-go-w-test/types"
	"github.com/alcb1310/bca-go-w-test/utils"
	"github.com/google/uuid"
)

func (s *Router) registerRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		type favContextKey string
		var c types.CreateCompany
		var ctx context.Context

		if r.Body == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"Message": "Invalid information",
			})
			return
		}

		if c.UserEmail != "" && !utils.IsValidEmail(c.UserEmail) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"Message": "Invalid information",
			})
			return
		}

		if c.ID == "" || c.Name == "" || c.UserEmail == "" || c.UserPassword == "" || c.UserName == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"Message": "Incomplete information",
			})
			return
		}

		k := favContextKey("company")
		ctx = context.WithValue(context.Background(), k, &c)
		tx, _ := s.db.Begin()
		defer tx.Rollback()

		sql := "INSERT INTO company (ruc, name, employees) VALUES ($1, $2, $3) RETURNING id, ruc, name, employees, is_active"
		var id, ruc, name string
		var employees uint8
		var is_active bool

		if err := tx.QueryRowContext(ctx, sql, c.ID, c.Name, c.Employees).Scan(&id, &ruc, &name, &employees, &is_active); err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		uuidId, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		co := types.Company{
			ID:        uuidId,
			Name:      name,
			RUC:       ruc,
			Employees: employees,
			IsActive:  is_active,
		}

		pass, err := utils.EncryptPasssword(c.UserPassword)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sql = "INSERT INTO \"user\" (email, name, password, company_id, role_id) VALUES ($1, $2, $3, $4, 'a')"
		if _, err := tx.ExecContext(ctx, sql, c.UserEmail, c.UserName, pass, co.ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tx.Commit()

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(co)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Router) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var lc types.LoginCredentials
		json.NewDecoder(r.Body).Decode(&lc)

		if lc.Email != "" && !utils.IsValidEmail(lc.Email) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"Message": "Invalid information",
			})
			return
		}

		if lc.Email == "" || lc.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var id, email, name, password, company_id, role_id string
		sql := "SELECT id, email, name, password, company_id, role_id from \"user\" where email = $1"
		if err := s.db.QueryRow(sql, lc.Email).Scan(&id, &email, &name, &password, &company_id, &role_id); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		isValid, _ := utils.ComparePassword(password, lc.Password)
		if !isValid {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		uId, _ := uuid.Parse(id)
		cId, _ := uuid.Parse(company_id)
		u := types.User{
			Id:        uId,
			Email:     email,
			Name:      name,
			CompanyId: cId,
			RoleId:    role_id,
		}

		json.NewEncoder(w).Encode(u)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
