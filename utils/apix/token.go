package apix

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"short-link/internal/consts"
	"short-link/logs"
)

func GetToken(id uint64, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   id,
		"username": username,
	})

	tokenString, err := token.SignedString([]byte(consts.TokenSecret))
	if err != nil {
		logs.Error(err, "GetToken failed")
	}
	return tokenString, err
}

func ParseToken(tokenString string) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(consts.TokenSecret), nil
	})

	if err != nil {
		return m, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return m, err
	}

	for k, v := range claims {
		m[k] = v
	}
	return m, nil
}
