package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterRoute(t *testing.T) {
	var s = NewRouter()
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

		// TODO: Make a valid request for tests to pass
		msg := map[string]string{
			"Message": "TEST",
		}
		msgJson, _ := json.Marshal(msg)
		msgReader := strings.NewReader(string(msgJson))
		req, _ = http.NewRequest(http.MethodPost, "/register", msgReader)
		res = httptest.NewRecorder()

		s.registerRoute(res, req)
		assertStatus(t, res.Code, http.StatusOK)
	})
}
