package middleware

import (
	"context"
	"dip/internal/models"
	"errors"
	"net/http"
	"strings"
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 48)),
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

func ValidateJWT(token string) (*TokenPayload, error) {
	t, err := jwt.ParseWithClaims(token, &TokenPayload{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(signature), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(*TokenPayload); ok && t.Valid && claims.ExpiresAt.After(time.Now()) {
		return claims, nil
	}

	return nil, errors.New("failed to decode token")
}

func AuthorizeToken() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			token := strings.Split(auth, " ")
			if len(token) < 2 {
				http.Error(w, http.ErrAbortHandler.Error(), http.StatusBadRequest)
				return
			}

			pl, err := ValidateJWT(token[1])
			ctx := r.Context()
			if err != nil {
				http.Error(w, http.ErrAbortHandler.Error(), http.StatusUnauthorized)
				return
			}
			ctx = context.WithValue(ctx, "creds", pl)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func ParseCredentials(ctx context.Context) (*TokenPayload, error) {
	tp, ok := ctx.Value("creds").(*TokenPayload)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return tp, nil
}
