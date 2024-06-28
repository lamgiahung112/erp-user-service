package handlers

import (
	"erp-user-service/middlewares"
	"erp-user-service/utils"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

type emailOtpLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Otp      string `json:"otp"`
}

func (hlr *HandlerConfig) LoginViaEmailOtp(w http.ResponseWriter, r *http.Request) {
	var requestPayload emailOtpLoginRequest
	err := hlr.readJSON(w, r, &requestPayload)
	if err != nil {
		hlr.errorJSON(w, hlr.ErrorFactory.Malformatted("request payload"))
		return
	}

	user, err := hlr.Models.Users.FindByEmail(requestPayload.Email)
	if err != nil {
		hlr.errorJSON(w, hlr.ErrorFactory.NotFound("user"), http.StatusNotFound)
		return
	}

	cachedOtp, err := hlr.Utils.Redis.GetUserLoginOtp(user.ID)
	if err != nil {
		hlr.errorJSON(w, err)
		return
	}

	sameOtp := strings.Compare(cachedOtp, requestPayload.Otp) == 0
	samePassword := user.PasswordMatches(requestPayload.Password)
	if !sameOtp || !samePassword {
		hlr.errorJSON(w, hlr.ErrorFactory.InvalidCredentials())
		return
	}

	refreshToken := uuid.NewString()
	token, err := hlr.Utils.Jwt.GenerateJwt(refreshToken, user.ToJwtUser().GetClaims())
	if err != nil {
		hlr.errorJSON(w, err)
		return
	}

	deviceInfo := r.Context().Value(middlewares.DeviceInfoMiddlewareResultKey).(*utils.DeviceInfo)
	deviceInfo.LoggedInAt = time.Now()
	err = hlr.Utils.Redis.StoreSessionInfo(user.ID, refreshToken, deviceInfo)
	if err != nil {
		hlr.errorJSON(w, hlr.ErrorFactory.StoreSessionFailed())
		return
	}

	err = hlr.Utils.Redis.RemoveUserLoginOtp(user.ID)
	if err != nil {
		hlr.errorJSON(w, err)
		return
	}

	resp := &jsonResponse{
		Error:   false,
		Message: "",
		Data: map[string]any{
			"token": token,
		},
	}
	hlr.writeJSON(w, http.StatusOK, resp)
}
