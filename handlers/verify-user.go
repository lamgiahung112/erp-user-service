package handlers

import (
	"erp-user-service/middlewares"
	"erp-user-service/utils"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (hlr *HandlerConfig) VerifyUser(w http.ResponseWriter, r *http.Request) {
	authClaims := r.Context().Value(middlewares.AuthenticationMiddlewareResultKey).(*jwt.MapClaims)
	user := hlr.Models.Users.ParseFromClaims(authClaims)
	refreshToken := fmt.Sprintf("%s", (*authClaims)["refresh"])

	deviceInfo := r.Context().Value(middlewares.DeviceInfoMiddlewareResultKey).(*utils.DeviceInfo)
	storedDeviceInfo, err := hlr.Utils.Redis.GetSessionInfo(user.ID, refreshToken)
	if err != nil {
		hlr.errorJSON(w, err)
		return
	}
	deviceInfo.LoggedInAt = storedDeviceInfo.LoggedInAt

	newRefreshToken := uuid.NewString()
	hlr.Utils.Redis.RemoveSessionInfo(user.ID, refreshToken)
	newJwt, err := hlr.Utils.Jwt.GenerateJwt(user.GetClaims(), newRefreshToken)

	if err != nil {
		hlr.errorJSON(w, err)
		return
	}

	err = hlr.Utils.Redis.StoreSessionInfo(user.ID, newRefreshToken, deviceInfo)
	if err != nil {
		hlr.errorJSON(w, err)
		return
	}

	resp := &jsonResponse{
		Error:   false,
		Message: "",
		Data: map[string]any{
			"token": newJwt,
			"user":  user,
		},
	}

	hlr.writeJSON(w, http.StatusOK, resp)
}
