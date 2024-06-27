package handlers

import (
	"erp-user-service/middlewares"
	"erp-user-service/utils"
	"erp-user-service/utils/rabbitmq"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (hlr *HandlerConfig) Login(w http.ResponseWriter, r *http.Request) {
	var requestPayload loginRequest
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

	// email OTP login
	if user.Is2FAEnabled && len(user.AuthenticatorSecretKey) == 0 {
		otp, err := hlr.Utils.Redis.StoreUserLoginOtp(user.ID)
		if err != nil {
			hlr.errorJSON(w, err)
			return
		}
		emailPayload := &rabbitmq.LoginOtpEmailPayload{
			ToAddress: user.Email,
			Otp:       otp,
			Title:     "One-time password to login to ERP",
			Username:  user.Name,
		}
		go hlr.Utils.EventEmitter.SendLoginOtpEmail(emailPayload)
		jsonResp := &jsonResponse{
			Error:   false,
			Message: "Please check your email to get the OTP to login",
		}
		hlr.writeJSON(w, http.StatusPartialContent, jsonResp)
		return
	}

	if user.Is2FAEnabled && len(user.AuthenticatorSecretKey) > 0 {
		jsonResp := &jsonResponse{
			Error:   false,
			Message: "You have 2FA enabled, please provide your 2FA secret code to login!",
		}
		hlr.writeJSON(w, http.StatusPartialContent, jsonResp)
		return
	}

	// Normal login
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

	resp := &jsonResponse{
		Error:   false,
		Message: "",
		Data: map[string]any{
			"token": token,
		},
	}
	hlr.writeJSON(w, http.StatusOK, resp)
}
