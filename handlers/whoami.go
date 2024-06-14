package handlers

import (
	"erp-user-service/middlewares"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func (hlr *HandlerConfig) WhoAmI(w http.ResponseWriter, r *http.Request) {
	authClaims := r.Context().Value(middlewares.AuthenticationMiddlewareResultKey).(*jwt.MapClaims)
	user := hlr.Models.Users.ParseFromClaims(authClaims)

	userFromDb, err := hlr.Models.Users.FindByUserID(user.ID)

	if err != nil {
		hlr.errorJSON(w, hlr.ErrorFactory.NotFound("user"))
		return
	}

	resp := &jsonResponse{
		Error:   false,
		Message: "",
		Data:    userFromDb,
	}
	hlr.writeJSON(w, http.StatusOK, resp)
}
