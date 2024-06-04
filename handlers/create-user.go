package handlers

import (
	"erp-user-service/data"
	"net/http"
)

type CreateUserRequestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Priority int16  `json:"priority"`
}

func (hlr *HandlerConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	var requestPayload CreateUserRequestPayload

	err := hlr.readJSON(w, r, &requestPayload)

	if err != nil {
		hlr.errorJSON(w, err)
		return
	}

	newUser := &data.Users{
		Email:    requestPayload.Email,
		Password: requestPayload.Password,
		Name:     requestPayload.Name,
		Priority: requestPayload.Priority,
	}

	err = hlr.Models.Users.Insert(newUser)

	if err != nil {
		hlr.errorJSON(w, err)
		return
	}

	jsonResp := jsonResponse{
		Error:   false,
		Message: "",
		Data:    newUser,
	}
	hlr.writeJSON(w, http.StatusAccepted, &jsonResp)
}
