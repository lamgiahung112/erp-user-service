package data

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtUtilities struct{}

func (user *Users) GetClaims() *map[string]any {
	return &map[string]any{
		"UserID": user.ID,
		"Email":  user.Email,
		"Name":   user.Name,
	}
}

func (*JwtUtilities) GenerateJwt(claimMap *map[string]any, refreshToken string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(jwtExpiration).Unix()
	claims["refresh"] = refreshToken

	for cl := range *claimMap {
		claims[cl] = (*claimMap)[cl]
	}

	signedToken, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (*JwtUtilities) VerifyJwt(token string) (*jwt.MapClaims, error) {
	parsedToken, err := jwt.NewParser().Parse(token, defaultKeyFunc)

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok || !parsedToken.Valid {
		return nil, errors.New("invalid credentials")
	}

	return &claims, nil
}

func defaultKeyFunc(t *jwt.Token) (interface{}, error) {
	return jwtKey, nil
}
