package handlers

import (
	"erp-user-service/middlewares"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func (hlr *HandlerConfig) GetDeviceList(w http.ResponseWriter, r *http.Request) {
	authClaims := r.Context().Value(middlewares.AuthenticationMiddlewareResultKey).(*jwt.MapClaims)
	user := hlr.Models.Users.ParseFromClaims(authClaims)

	devices, err := hlr.Utils.Redis.GetAllSessionsOfUser(user.ID)
	if err != nil {
		hlr.errorJSON(w, err)
		return
	}

	resp := &jsonResponse{
		Error:   false,
		Message: "",
		Data:    devices,
	}

	hlr.writeJSON(w, http.StatusOK, resp)
}
