package middleware

import (
	"dip/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const signature = "7fec87f134b063cd0546d7059f7d1acb4c365229b9dd4c66259c67b65ee33a65"

type TokenPayload struct {
	Id           int
	Login        string
	Phone        string
	Email        string
	IsMaintainer bool
	jwt.RegisteredClaims
}

func GenerateJWT(p *models.Person) (string, error) {
	tp := &TokenPayload{
		p.Id, p.Name, p.Phone, p.Email, p.IsMaintainer,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 20)),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tp)
	st, err := token.SignedString([]byte(signature))
	if err != nil {
		return "", err
	}

	return st, nil
}

func ValidateJWT(token string) (TokenPayload, error) {
	t, _ := jwt.ParseWithClaims(token, &TokenPayload{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(signature), nil
	})

	if claims, ok := t.Claims.(*TokenPayload); ok && t.Valid {
		return *claims, nil
	}

	return TokenPayload{}, nil
}
