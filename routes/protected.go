package routes

import (
	"database/sql"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

type ProtectedRouter struct {
	*mux.Router

	db *sql.DB
}

func (s *Router) protectedRoutes() {
	p := &ProtectedRouter{
		Router: s.PathPrefix("/api/v1").Subrouter(),
		db:     s.db,
	}
	s.Use(s.jsonResponse)
	p.Use(p.authVerify)

	p.HandleFunc("/logout", p.handleLogout).Methods(http.MethodGet)
	p.HandleFunc("/", p.handleBCAHome).Methods(http.MethodGet)
}

func (s *ProtectedRouter) handleBCAHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/bca/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusTeapot)
		return
	}
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)

	return
}

func (s *ProtectedRouter) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ctxPayload, _ := getMyPaload(r)
		sql := "DELETE FROM logged_in_user WHERE user_id = $1"
		if _, err := s.db.Exec(sql, ctxPayload.Id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}
