package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alcb1310/bca-go-w-test/utils"
	"github.com/google/uuid"
)

func (s *Router) jsonResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (s *ProtectedRouter) authVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("bca")

		token := strings.Split(cookie.String(), "=")
		if len(token) != 2 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		secretKey := os.Getenv("SECRET")
		maker, err := utils.NewJWTMaker(secretKey)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tokenData, err := maker.VerifyToken(token[1])
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		marshalStr, _ := json.Marshal(tokenData)
		ctx := r.Context()
		ctx = context.WithValue(r.Context(), "token", marshalStr)
		r = r.Clone(ctx)

		var t string
		sql := "SELECT token FROM logged_in_user WHERE user_id = $1"
		if err := s.db.QueryRow(sql, tokenData.ID).Scan(&t); err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if !bytes.Equal([]byte(token[1]), []byte(t)) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if r.URL.Path == "/api/v1/edit-user" {
			if strings.TrimRight(tokenData.Role, " ") != "a" {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

type contextPayload struct {
	Id         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	CompanyId  uuid.UUID `json:"company_id"`
	Role       string    `json:"role"`
	IsLoggedIn bool      `json:"is_logged_in"`
	IssuedAt   time.Time `json:"issued_at"`
	ExpiredAt  time.Time `json:"expired_at"`
}

func getMyPaload(r *http.Request) (contextPayload, error) {
	ctx := r.Context()
	val := ctx.Value("token")

	x, ok := val.([]byte)
	if !ok {
		return contextPayload{}, errors.New("Unable to load context")
	}
	var p contextPayload
	err := json.Unmarshal(x, &p)
	if err != nil {
		return contextPayload{}, errors.New("Unable to parse context")
	}
	return p, nil
}
