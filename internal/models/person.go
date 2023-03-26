package models

import (
	"errors"
	"net/mail"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type Person struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Settings string `json:"settings"`
	Phone    string `json:"phones"`
}

func (p *Person) Validate() error {
	match, _ := regexp.MatchString(`^[\+]?[(]?[0-9]{3}[)]?[-\s\.]?[0-9]{3}[-\s\.]?[0-9]{4,6}$`, p.Phone)
	if !match {
		return errors.New("Incorrect phone number format")
	}
	_, err := mail.ParseAddress(p.Email)
	if err != nil {
		return errors.New("Incorrect email format")
	}
	b, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	p.Password = string(b)
	return nil
}
