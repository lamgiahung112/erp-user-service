package handlers

import (
	"erp-user-service/middlewares"
	"erp-user-service/utils"
	"log"
	"net/http"

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
	}

	user, err := hlr.Models.Users.FindByEmail(requestPayload.Email)
	if err != nil {
		hlr.errorJSON(w, hlr.ErrorFactory.NotFound("user"), http.StatusNotFound)
	}

	refreshToken := uuid.NewString()
	token, err := hlr.Utils.Jwt.GenerateJwt(user.ToJwtUser().GetClaims(), refreshToken)
	if err != nil {
		hlr.errorJSON(w, err)
		return
	}

	iplocation := r.Context().Value(middlewares.IpLocationMiddlewareResultKey).(*utils.IpLocationData)
	if iplocation == nil {
		hlr.errorJSON(w, hlr.ErrorFactory.Unexpected())
		return
	}

	userAgent := r.Header.Get("user-agent")
	deviceInfo, err := hlr.Utils.DeviceInfo.GetDevice(userAgent, iplocation)
	if err != nil {
		hlr.errorJSON(w, err)
		return
	}
	log.Println(&deviceInfo)
	err = hlr.Utils.Redis.StoreSessionInfo(user.ID, refreshToken, deviceInfo)
	if err != nil {
		hlr.errorJSON(w, hlr.ErrorFactory.StoreSessionFailed())
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
