package main

import (
	"erp-user-service/data"
	"net/http"
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

	id, err := app.Models.Users.Insert(newUser)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	jsonResp := jsonResponse{
		Error:   false,
		Message: "",
		Data:    id,
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

	jwt, err := user.GenerateJwt()

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "",
		Data: map[string]string{
			"token": jwt,
		},
	}
	app.writeJSON(w, http.StatusAccepted, &resp)
}
