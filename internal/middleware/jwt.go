package middleware

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var signature = []byte("secretKeytoSign")

type TokenPayload struct {
	Login string
	Phone string
	Email string
	jwt.RegisteredClaims
}

func GenerateJWT(phone, email, login string) (string, error) {
	tp := &TokenPayload{
		login, phone, email,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 20)),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tp)
	st, err := token.SignedString(signature)
	if err != nil {
		return "", err
	}

	return st, nil
}

func ValidateJWT(token string) (TokenPayload, error) {
	t, _ := jwt.ParseWithClaims(token, &TokenPayload{}, func(t *jwt.Token) (interface{}, error) {
		return signature, nil
	})

	if claims, ok := t.Claims.(*TokenPayload); ok && t.Valid {
		return *claims, nil
	}

	return TokenPayload{}, nil
}
