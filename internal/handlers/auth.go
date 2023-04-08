package handlers

import (
	"dip/internal/middleware"
	"dip/internal/models"
	"dip/internal/store"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func LogIn(store store.Store, w http.ResponseWriter, r *http.Request) (int, string, error) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	req := &request{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return http.StatusBadRequest, "", err
	}
	p, err := store.Person().GetByEmail(req.Email)
	if err != nil {
		return http.StatusUnauthorized, "", err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(req.Password)); err != nil {
		return http.StatusUnauthorized, "", err
	}

	token, err := middleware.GenerateJWT(p.Phone, p.Email, p.Name)
	if err != nil {
		return http.StatusBadRequest, "", err
	}

	return http.StatusOK, token, nil
}

func SignUp(store store.Store, w http.ResponseWriter, r *http.Request) (error, int) {
	type request struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}

	req := &request{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err, http.StatusBadRequest
	}
	p := &models.Person{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Phone:    req.Phone,
		Settings: "",
	}

	if err := store.Person().Create(p); err != nil {
		return err, http.StatusUnprocessableEntity
	}

	return nil, http.StatusCreated
}
