package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alcb1310/bca-go-w-test/types"
	"github.com/gorilla/mux"
)

func TestRegisterRoute(t *testing.T) {
	var s = &Router{
		Router: mux.NewRouter(),
	}
	t.Run("it should return 405 for all methods but POST to the /register route", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/register", nil)
		res := httptest.NewRecorder()

		s.registerRoute(res, req)
		assertStatus(t, res.Code, http.StatusMethodNotAllowed)

		req, _ = http.NewRequest(http.MethodPut, "/register", nil)
		res = httptest.NewRecorder()

		s.registerRoute(res, req)
		assertStatus(t, res.Code, http.StatusMethodNotAllowed)

		req, _ = http.NewRequest(http.MethodPatch, "/register", nil)
		res = httptest.NewRecorder()

		s.registerRoute(res, req)
		assertStatus(t, res.Code, http.StatusMethodNotAllowed)

		req, _ = http.NewRequest(http.MethodDelete, "/register", nil)
		res = httptest.NewRecorder()

		s.registerRoute(res, req)
		assertStatus(t, res.Code, http.StatusMethodNotAllowed)

		req, _ = http.NewRequest(http.MethodPost, "/register", nil)
		res = httptest.NewRecorder()

		s.registerRoute(res, req)
		assertStatus(t, res.Code, http.StatusBadRequest)
	})

	t.Run("it should throw 400 when incomplete/incorrect data is sent", func(t *testing.T) {
		msg := map[string]string{
			"Message": "TEST",
		}
		msgJson, _ := json.Marshal(msg)
		msgReader := strings.NewReader(string(msgJson))
		req, _ := http.NewRequest(http.MethodPost, "/register", msgReader)
		res := httptest.NewRecorder()

		s.registerRoute(res, req)
		assertStatus(t, res.Code, http.StatusBadRequest)
		json.Unmarshal([]byte(res.Body.String()), &msg)
		assertMessage(t, msg["Message"], "Incomplete information")

		msg = map[string]string{
			"id": "123",
		}
		msgJson, _ = json.Marshal(msg)
		msgReader = strings.NewReader(string(msgJson))
		req, _ = http.NewRequest(http.MethodPost, "/register", msgReader)
		res = httptest.NewRecorder()
		s.registerRoute(res, req)
		assertStatus(t, res.Code, http.StatusBadRequest)
		json.Unmarshal([]byte(res.Body.String()), &msg)
		assertMessage(t, msg["Message"], "Incomplete information")

		msg = map[string]string{
			"id":        "123",
			"name":      "test company",
			"employees": "lj2",
		}
		msgJson, _ = json.Marshal(msg)
		msgReader = strings.NewReader(string(msgJson))
		req, _ = http.NewRequest(http.MethodPost, "/register", msgReader)
		res = httptest.NewRecorder()
		s.registerRoute(res, req)
		assertStatus(t, res.Code, http.StatusBadRequest)
		json.Unmarshal([]byte(res.Body.String()), &msg)
		assertMessage(t, msg["Message"], "Invalid information")

		msg = map[string]string{
			"id":        "123",
			"name":      "test company",
			"employees": "256",
		}
		msgJson, _ = json.Marshal(msg)
		msgReader = strings.NewReader(string(msgJson))
		req, _ = http.NewRequest(http.MethodPost, "/register", msgReader)
		res = httptest.NewRecorder()
		s.registerRoute(res, req)
		assertStatus(t, res.Code, http.StatusBadRequest)
		json.Unmarshal([]byte(res.Body.String()), &msg)
		assertMessage(t, msg["Message"], "Invalid information")

		msg = map[string]string{
			"id":   "123",
			"name": "test company",
		}
		msgJson, _ = json.Marshal(msg)
		msgReader = strings.NewReader(string(msgJson))
		req, _ = http.NewRequest(http.MethodPost, "/register", msgReader)
		res = httptest.NewRecorder()
		s.registerRoute(res, req)
		assertStatus(t, res.Code, http.StatusBadRequest)
		json.Unmarshal([]byte(res.Body.String()), &msg)
		assertMessage(t, msg["Message"], "Incomplete information")

		msg = map[string]string{
			"id":         "123",
			"name":       "test company",
			"user_email": "invalid",
		}
		msgJson, _ = json.Marshal(msg)
		msgReader = strings.NewReader(string(msgJson))
		req, _ = http.NewRequest(http.MethodPost, "/register", msgReader)
		res = httptest.NewRecorder()
		s.registerRoute(res, req)
		assertStatus(t, res.Code, http.StatusBadRequest)
		json.Unmarshal([]byte(res.Body.String()), &msg)
		assertMessage(t, msg["Message"], "Invalid information")

		msg = map[string]string{
			"id":         "123",
			"name":       "test company",
			"user_email": "valid@email.com",
		}
		msgJson, _ = json.Marshal(msg)
		msgReader = strings.NewReader(string(msgJson))
		req, _ = http.NewRequest(http.MethodPost, "/register", msgReader)
		res = httptest.NewRecorder()
		s.registerRoute(res, req)
		assertStatus(t, res.Code, http.StatusBadRequest)
		json.Unmarshal([]byte(res.Body.String()), &msg)
		assertMessage(t, msg["Message"], "Incomplete information")

		msg = map[string]string{
			"id":            "123",
			"name":          "test company",
			"user_email":    "valid@email.com",
			"user_password": "password123",
		}
		msgJson, _ = json.Marshal(msg)
		msgReader = strings.NewReader(string(msgJson))
		req, _ = http.NewRequest(http.MethodPost, "/register", msgReader)
		res = httptest.NewRecorder()
		s.registerRoute(res, req)
		assertStatus(t, res.Code, http.StatusBadRequest)
		json.Unmarshal([]byte(res.Body.String()), &msg)
		assertMessage(t, msg["Message"], "Incomplete information")
	})

	t.Run("on correct data it should success", func(t *testing.T) {
		c := types.CreateCompany{
			ID:           "123",
			Name:         "Test Company",
			UserEmail:    "valid@email.com",
			UserName:     "Test User",
			UserPassword: "password123",
			Employees:    2,
		}

		msgJson, _ := json.Marshal(c)
		msgReader := strings.NewReader(string(msgJson))
		req, _ := http.NewRequest(http.MethodPost, "/register", msgReader)
		res := httptest.NewRecorder()
		s.registerRoute(res, req)
		assertStatus(t, res.Code, http.StatusOK)

	})
}

func assertMessage(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("incorrect message, got %s, want %s", got, want)
	}
}
