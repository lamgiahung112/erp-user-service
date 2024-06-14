package handlers

import (
	"erp-user-service/middlewares"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func (hlr *HandlerConfig) AddThirdPartyAuthenticator(w http.ResponseWriter, r *http.Request) {
	authClaims := r.Context().Value(middlewares.AuthenticationMiddlewareResultKey).(*jwt.MapClaims)
	user := hlr.Models.Users.ParseFromClaims(authClaims)

	qrData, secretKey, err := hlr.Utils.QR.GenerateQrCodeData(user.Email)
	if err != nil {
		hlr.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	dbUser, err := hlr.Models.Users.FindByUserID(user.ID)
	if err != nil {
		hlr.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	if dbUser.Is2FAEnabled == false {
		hlr.errorJSON(w, hlr.ErrorFactory.NotEnabled("2FA"), http.StatusForbidden)
		return
	}
	if len(dbUser.AuthenticatorSecretKey) > 0 {
		hlr.errorJSON(w, hlr.ErrorFactory.AlreadyExists("an authentication key"), http.StatusInternalServerError)
		return
	}

	dbUser.AuthenticatorSecretKey = secretKey
	err = hlr.Models.Users.Save(dbUser)
	if err != nil {
		hlr.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	hlr.writeImage(w, qrData)
}
