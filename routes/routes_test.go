package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestMainRoute(t *testing.T) {
	s := &Router{
		Router: mux.NewRouter(),
	}
	t.Run("should return 200 on GET request", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		s.mainRoute(res, req)
		assertStatus(t, res.Code, http.StatusOK)

		var x map[string]string

		json.Unmarshal([]byte(res.Body.String()), &x)

		if x["message"] != "Server Ok" {
			t.Error("Incorrect message received")
		}
	})

	t.Run("should return 405 on all other methods", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/", nil)
		res := httptest.NewRecorder()
		s.mainRoute(res, req)
		assertStatus(t, res.Code, http.StatusMethodNotAllowed)

		req, _ = http.NewRequest(http.MethodPatch, "/", nil)
		res = httptest.NewRecorder()
		s.mainRoute(res, req)
		assertStatus(t, res.Code, http.StatusMethodNotAllowed)

		req, _ = http.NewRequest(http.MethodPut, "/", nil)
		res = httptest.NewRecorder()
		s.mainRoute(res, req)
		assertStatus(t, res.Code, http.StatusMethodNotAllowed)

		req, _ = http.NewRequest(http.MethodDelete, "/", nil)
		res = httptest.NewRecorder()
		s.mainRoute(res, req)
		assertStatus(t, res.Code, http.StatusMethodNotAllowed)
	})
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}
