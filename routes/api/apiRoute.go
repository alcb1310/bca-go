package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/alcb1310/bca-go-w-test/database"
	"github.com/gorilla/mux"
)

type router struct {
	*mux.Router
	db *database.Database
}

func NewRouter() *router {
	ro := mux.NewRouter()
	conn, err := database.ConnectDB()
	if err != nil {
		log.Panic(fmt.Sprintf("Unable to connect to the database, error: %s", err.Error()))
	}
	r := &router{
		Router: ro.PathPrefix("/api/v1").Subrouter(),
		db:     conn,
	}

	r.routes()
	return r
}

func (s *router) routes() {
	s.HandleFunc("/", s.mainRoute).Methods(http.MethodGet)
}

func (s *router) mainRoute(w http.ResponseWriter, r *http.Request) {
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
}
