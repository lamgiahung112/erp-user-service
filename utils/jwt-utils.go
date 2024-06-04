package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtUtilities struct {
	expirationPeriod time.Duration
	key              []byte
}

func (config *JwtUtilities) GenerateJwt(userID string, refreshToken string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(config.expirationPeriod).Unix()
	claims["refresh"] = refreshToken
	claims["userID"] = userID

	signedToken, err := token.SignedString(config.key)

	if err != nil {
		return "", ErrorFactory.Unexpected()
	}

	return signedToken, nil
}

func (config *JwtUtilities) VerifyJwt(token string) (*jwt.MapClaims, error) {
	parsedToken, err := jwt.NewParser().Parse(token, config.defaultKeyFunc)

	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok || !parsedToken.Valid {
		return nil, ErrorFactory.InvalidCredentials()
	}

	return &claims, nil
}

func (config *JwtUtilities) defaultKeyFunc(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, ErrorFactory.Malformatted("credentials")
	}

	return config.key, nil
}
