package routes

import (
	"net/http"
	"text/template"

	"github.com/alcb1310/bca-go-w-test/database"
	"github.com/alcb1310/bca-go-w-test/utils"
	"github.com/gorilla/mux"
)

type ProtectedRouter struct {
	*mux.Router

	db *database.Database
}

func (s *Router) protectedRoutes() {
	p := &ProtectedRouter{
		Router: s.PathPrefix(utils.HTML_PATH_PREFIX).Subrouter(),
		db:     s.db,
	}
	p.Use(p.authVerify)

	p.HandleFunc("/logout", p.handleLogout)
	p.HandleFunc("/", p.handleBCAHome)

	p.HandleFunc("/usuarios", p.handleUsers)
	p.transactionsRoutes()
}
func (s *ProtectedRouter) handleBCAHome(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)
	switch r.Method {
	case http.MethodGet:
		tmpl, err := template.ParseFiles(utils.TEMPLATE_DIR + "/bca/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusTeapot)
			return
		}
		type Ret struct {
			UserName string
			Title    string
			Links    utils.LinksType
		}
		retData := Ret{
			UserName: ctxPayload.Name,
			Title:    "BCA - home",
			Links:    *utils.Links,
		}

		tmpl.ExecuteTemplate(w, "index.html", retData)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *ProtectedRouter) handleLogout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ctxPayload, _ := getMyPaload(r)
		if err := s.db.Logout(ctxPayload.Id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
