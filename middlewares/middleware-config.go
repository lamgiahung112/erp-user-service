package middlewares

import (
	"encoding/json"
	"erp-user-service/data"
	"erp-user-service/factory"
	"erp-user-service/utils"
	"net/http"
)

type MiddlewareConfig struct {
	Utils        *utils.AppUtilities
	ErrorFactory *factory.ErrorFactory
	Models       *data.Models
}

type key string

func (c key) String() string {
	return string(c)
}

const (
	AuthenticationMiddlewareResultKey = key("auth-result-key")
	IpLocationMiddlewareResultKey     = key("iplocation-result-key")
)

func New() *MiddlewareConfig {
	return &MiddlewareConfig{
		Utils:        utils.New(),
		ErrorFactory: &factory.ErrorFactory{},
		Models:       data.New(),
	}
}

func (md *MiddlewareConfig) serveError(err error, w http.ResponseWriter, status int) {
	errResp := map[string]any{
		"error":   true,
		"message": err.Error(),
	}

	out, _ := json.Marshal(&errResp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(out)
}
