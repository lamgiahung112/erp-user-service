package main

import (
	"erp-user-service/data"
	"errors"
	"fmt"
	"net/http"
	"time"

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

	jwtClaims := user.ToJwtUser().GetClaims()
	refreshToken := uuid.NewString()

	jwt, err := app.Utils.Jwt.GenerateJwt(jwtClaims, refreshToken)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	iplocation, err := app.getClientIpLocation(r)

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
			"token": jwt,
			"user":  user.ToJwtUser(),
		},
	}
	app.writeJSON(w, http.StatusAccepted, &resp)
}

func (app *Config) Verify(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")[7:]
	claims, err := app.Utils.Jwt.VerifyJwt(token)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	}

	refreshToken := fmt.Sprintf("%s", (*claims)["refresh"])
	user := app.Models.Users.ParseFromClaims(claims)

	if len(refreshToken) == 0 || len(user.Email) == 0 || len(user.ID) == 0 || len(user.Name) == 0 {
		app.errorJSON(w, app.ErrorFactory.InvalidCredentials())
		return
	}

	storedDeviceInfo, err := app.Utils.Redis.GetSessionInfo(user.ID, refreshToken)
	println(storedDeviceInfo)
	if err != nil {
		app.errorJSON(w, app.ErrorFactory.InvalidCredentials())
		return
	}

	iplocation, err := app.getClientIpLocation(r)
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

	isDeviceDifferent := storedDeviceInfo.Compare(deviceInfo)
	if isDeviceDifferent {
		app.errorJSON(w, app.ErrorFactory.InvalidCredentials())
		return
	}
	deviceInfo.LoginPortal = storedDeviceInfo.LoginPortal
	deviceInfo.LoggedInAt = time.Now()

	err = app.Utils.Redis.RemoveRefreshToken(user.ID, refreshToken)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	newRefreshToken := uuid.NewString()
	err = app.Utils.Redis.StoreSessionInfo(user.ID, newRefreshToken, deviceInfo)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	jwt, err := app.Utils.Jwt.GenerateJwt(user.GetClaims(), newRefreshToken)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "",
		Data: map[string]any{
			"token": jwt,
			"user":  user,
		},
	}

	app.writeJSON(w, http.StatusAccepted, &resp)
}

func (app *Config) TestGRPC(w http.ResponseWriter, r *http.Request) {
	go app.LogViaGRPC("Test grpc", "Testing very carefully...")
	app.writeJSON(w, http.StatusAccepted, map[string]any{})
}
