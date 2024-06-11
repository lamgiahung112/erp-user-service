package handlers

import (
	"erp-user-service/middlewares"
	"erp-user-service/utils"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// RefreshUserSession verifies user provided by the value of middlewares.AuthenticationMiddlewareResultKey and returns a new jwt
//
// Steps:
//  1. Parse claims from Authentication middleware to get user and refresh token
//  2. Get device info from DeviceInfo middleware
//  3. Get the cached device info
//  4. Compare the cached device and current device
//  5. Create new refresh token, new jwt
//  6. Invalidate refresh token from Authentication middleware
//  7. Store new refresh token along with new device info
//  8. Send to user the new jwt
func (hlr *HandlerConfig) RefreshUserSession(w http.ResponseWriter, r *http.Request) {
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
	err = hlr.Utils.Redis.RemoveSessionInfo(user.ID, refreshToken)
	if err != nil {
		hlr.errorJSON(w, err)
		return
	}
	newJwt, err := hlr.Utils.Jwt.GenerateJwt(newRefreshToken, user.GetClaims())

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
		},
	}

	hlr.writeJSON(w, http.StatusOK, resp)
}
