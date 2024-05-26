package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtUtilities struct {
	expirationPeriod time.Duration
	key              []byte
}

func (config *JwtUtilities) GenerateJwt(claimMap *map[string]any, refreshToken string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(config.expirationPeriod).Unix()
	claims["refresh"] = refreshToken

	for cl := range *claimMap {
		claims[cl] = (*claimMap)[cl]
	}

	signedToken, err := token.SignedString(config.key)

	if err != nil {
		return "", err
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
		return nil, errors.New("invalid credentials")
	}

	return &claims, nil
}

func (config *JwtUtilities) defaultKeyFunc(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("unexpected signing method")
	}

	return config.key, nil
}
