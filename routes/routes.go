package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/alcb1310/bca-go-w-test/database"
	"github.com/gorilla/mux"
)

type Router struct {
	*mux.Router

	db *sql.DB
}

func NewRouter() *Router {
	conn, err := database.ConnectDB()
	if err != nil {
		log.Panic(fmt.Sprintf("Unable to connect to the database, error: %s", err.Error()))
	}
	r := &Router{
		Router: mux.NewRouter(),
		db:     conn,
	}

	r.routes()
	return r
}

func (s *Router) routes() {
	s.HandleFunc("/", s.mainRoute)
	s.HandleFunc("/register", s.registerRoute)
	s.HandleFunc("/login", s.handleLogin)

	s.protectedRoutes()
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
