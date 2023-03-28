package handlers

import (
	"dip/internal/models"
	"dip/internal/store"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func LogIn(store store.Store, w http.ResponseWriter, r *http.Request) (error, int) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	req := &request{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err, http.StatusBadRequest
	}
	p, err := store.Person().GetByEmail(req.Email)
	if err != nil {
		return err, http.StatusUnauthorized
	}
	if err = bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(req.Password)); err != nil {
		return err, http.StatusUnauthorized
	}
	return nil, http.StatusOK
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
