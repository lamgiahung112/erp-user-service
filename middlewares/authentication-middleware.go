package middlewares

import (
	"context"
	"erp-user-service/utils"
	"fmt"
	"net/http"
	"strings"
)

func (md *MiddlewareConfig) Authenticated(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
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
		if user == nil || len(refreshToken) == 0 || len(user.Email) == 0 || len(user.ID) == 0 || len(user.Name) == 0 {
			md.serveError(md.ErrorFactory.InvalidCredentials(), w, http.StatusUnauthorized)
			return
		}

		deviceInfo := r.Context().Value(DeviceInfoMiddlewareResultKey).(*utils.DeviceInfo)
		storedDeviceInfo, err := md.Utils.Redis.GetSessionInfo(user.ID, refreshToken)

		if err != nil {
			md.serveError(err, w, http.StatusUnauthorized)
			return
		}

		isDifferentDevice := deviceInfo.Compare(storedDeviceInfo)

		if isDifferentDevice {
			md.serveError(md.ErrorFactory.NotFound("device"), w, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), AuthenticationMiddlewareResultKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
