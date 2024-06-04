package handlers

import (
	"erp-user-service/middlewares"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func (hlr *HandlerConfig) RevokeUserDevices(w http.ResponseWriter, r *http.Request) {
	authClaims := r.Context().Value(middlewares.AuthenticationMiddlewareResultKey).(*jwt.MapClaims)
	user := hlr.Models.Users.ParseFromClaims(authClaims)

	err := hlr.Utils.Redis.RevokeAllUserSessions(user.ID)
	if err != nil {
		hlr.errorJSON(w, err)
		return
	}

	resp := &jsonResponse{
		Error:   false,
		Message: "All sessions have been revoked",
		Data:    nil,
	}

	hlr.writeJSON(w, http.StatusOK, resp)
}
