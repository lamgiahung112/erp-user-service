package handlers

import (
	"erp-user-service/middlewares"
	"erp-user-service/utils"
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

	refreshToken := uuid.NewString()
	token, err := hlr.Utils.Jwt.GenerateJwt(user.ID, refreshToken)
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
			"user":  user.ToJwtUser(),
		},
	}
	hlr.writeJSON(w, http.StatusOK, resp)
}
