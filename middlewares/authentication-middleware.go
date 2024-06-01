package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

func (md *MiddlewareConfig) Authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) < 7 || !strings.HasPrefix(authHeader, "Bearer ") {
			md.serveError(md.ErrorFactory.InvalidCredentials(), w, http.StatusUnauthorized)
			return
		}

		token := authHeader[7:]
		claims, err := md.Utils.Jwt.VerifyJwt(token)
		if err != nil {
			md.serveError(err, w, http.StatusUnauthorized)
			return
		}

		refreshToken := fmt.Sprintf("%s", (*claims)["refresh"])
		user := md.Models.Users.ParseFromClaims(claims)

		if len(refreshToken) == 0 || len(user.Email) == 0 || len(user.ID) == 0 || len(user.Name) == 0 {
			md.serveError(md.ErrorFactory.InvalidCredentials(), w, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), AuthenticationMiddlewareResultKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
