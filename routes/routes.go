package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	*mux.Router
}

func NewRouter() *Router {
	r := &Router{
		Router: mux.NewRouter(),
	}

	r.routes()
	return r
}

func (s *Router) routes() {
	s.HandleFunc("/", s.mainRoute)
}

func (s *Router) mainRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		msg := map[string]string{
			"message": "Server Ok",
		}
		msgJson, err := json.Marshal(msg)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(msgJson)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}
