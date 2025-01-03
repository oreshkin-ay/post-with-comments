package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/oreshkin/posts/internal/pkg/jwt"
	"github.com/oreshkin/posts/internal/users"
)

var UserCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			fields := strings.Fields(header)
			if len(fields) != 2 || fields[0] != "Bearer" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)

				if err := json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token"}); err != nil {
					http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
					return
				}
				return
			}

			tokenStr := fields[1]

			username, err := jwt.ParseToken(tokenStr)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)

				if err := json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token"}); err != nil {
					http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
					return
				}
				return
			}

			id, err := users.GetUserIdByUsername(username)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			user := users.User{Username: username, ID: strconv.Itoa(id)}
			ctx := context.WithValue(r.Context(), UserCtxKey, &user)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func ForContext(ctx context.Context) *users.User {
	raw, _ := ctx.Value(UserCtxKey).(*users.User)
	return raw
}
