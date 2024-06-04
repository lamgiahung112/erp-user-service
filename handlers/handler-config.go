package handlers

import (
	"encoding/json"
	"erp-user-service/data"
	"erp-user-service/factory"
	"erp-user-service/utils"
	"errors"
	"io"
	"net/http"
)

type HandlerConfig struct {
	Models       *data.Models
	Utils        *utils.AppUtilities
	ErrorFactory *factory.ErrorFactory
}

func New() *HandlerConfig {
	return &HandlerConfig{
		Models:       data.New(),
		ErrorFactory: &factory.ErrorFactory{},
		Utils:        utils.New(),
	}
}

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (hlr *HandlerConfig) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1408576 // 1MB

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)

	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})

	if err != io.EOF {
		return errors.New("body must only have only 1 JSON value")
	}

	return nil
}

func (hlr *HandlerConfig) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)

	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)

	if err != nil {
		return err
	}

	return nil
}

func (hlr *HandlerConfig) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	payload := jsonResponse{
		Error:   true,
		Message: err.Error(),
	}

	return hlr.writeJSON(w, statusCode, payload)
}
