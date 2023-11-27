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
	s.HandleFunc("/login", s.handleLogin).Methods(http.MethodPost)
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

func (s *router) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		type loginPayload struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		var lp loginPayload
		err := json.NewDecoder(r.Body).Decode(&lp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if lp.Email == "" || lp.Password == "" {
			msg := map[string]string{
				"message": "Información Incompleta",
			}
			msgJson, _ := json.Marshal(msg)
			http.Error(w, string(msgJson), http.StatusBadRequest)
			return
		}

		token, err := s.db.Login(lp.Email, lp.Password)
		if err != nil {
			msg := map[string]string{
				"message": err.Error(),
			}
			msgJson, _ := json.Marshal(msg)
			http.Error(w, string(msgJson), http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		msg := map[string]string{
			"token": token,
		}

		msgJson, _ := json.Marshal(msg)
		w.Write(msgJson)
		return
	}
}
