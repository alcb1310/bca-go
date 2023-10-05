package routes

import (
	"encoding/json"
	"net/http"

	"github.com/alcb1310/bca-go-w-test/types"
)

func (s *Router) registerRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var c types.CreateCompany

		if r.Body == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}
