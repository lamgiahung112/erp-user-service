package main

import (
	"erp-user-service/data"
	"net/http"

	"github.com/google/uuid"
)

type CreateUserRequestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) CreateUser(w http.ResponseWriter, r *http.Request) {
	var requestPayload CreateUserRequestPayload

	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	newUser := data.Users{
		Email:    requestPayload.Email,
		Password: requestPayload.Password,
		Name:     requestPayload.Name,
	}

	err = app.Models.Users.Insert(&newUser)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	jsonResp := jsonResponse{
		Error:   false,
		Message: "",
		Data:    &newUser,
	}
	app.writeJSON(w, http.StatusAccepted, &jsonResp)
}

func (app *Config) Login(w http.ResponseWriter, r *http.Request) {
	var requestPayload LoginRequestPayload
	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	user, err := app.Models.Users.FindByEmail(requestPayload.Email)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	matchesPassword := user.PasswordMatches(requestPayload.Password)

	if !matchesPassword {
		app.errorJSON(w, err)
		return
	}

	jwtClaims := user.GetClaims()
	refreshToken := uuid.NewString()

	jwt, err := app.Utils.Jwt.GenerateJwt(jwtClaims, refreshToken)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	iplocation, err := app.getClientDeviceInfo(r)

	if err != nil {
		app.errorJSON(w, err)
		return
	}
	userAgent := app.getUserAgent(r)
	deviceInfo, err := app.Utils.DeviceInfo.GetDevice(userAgent, iplocation)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	deviceInfo.LoginPortal = "Local"
	app.Utils.Redis.StoreSessionInfo(user.ID, refreshToken, deviceInfo)

	resp := jsonResponse{
		Error:   false,
		Message: "",
		Data: map[string]any{
			"token":      jwt,
			"deviceInfo": deviceInfo,
		},
	}
	app.writeJSON(w, http.StatusAccepted, &resp)
}

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	// bearerToken := r.Header.Get("Authorization")
	// token := bearerToken[7:len(bearerToken)]

}
