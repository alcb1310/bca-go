package routes

import (
	"net/http"
	"text/template"

	"github.com/alcb1310/bca-go-w-test/database"
	"github.com/gorilla/mux"
)

type ProtectedRouter struct {
	*mux.Router

	db *database.Database
}

func (s *Router) protectedRoutes() {
	p := &ProtectedRouter{
		Router: s.PathPrefix("/api/v1").Subrouter(),
		db:     s.db,
	}
	s.Use(s.jsonResponse)
	p.Use(p.authVerify)

	p.HandleFunc("/logout", p.handleLogout).Methods(http.MethodGet)
	p.HandleFunc("/change-password", p.tmplChangePassword).Methods(http.MethodGet)
	p.HandleFunc("/edit-user", p.handleEditUser)
	p.HandleFunc("/", p.handleBCAHome).Methods(http.MethodGet)

	p.HandleFunc("/users", p.handleUsers)
	p.HandleFunc("/users/add", p.tmplAddUser)
	p.HandleFunc("/users/{userId}", p.handleSimpleUser)
}

func (s *ProtectedRouter) handleEditUser(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/bca/edit-user.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusTeapot)
		return
	}
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)

	return
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
		if err := s.db.Logout(ctxPayload.Id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}
