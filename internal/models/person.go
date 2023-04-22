package models

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"golang.org/x/crypto/bcrypt"
)

type Person struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	Settings     string `json:"settings"`
	Phone        string `json:"phones"`
	IsMaintainer bool   `json:"isMaintainer"`
}

type PersonInTask struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func (p *Person) Validate() error {
	if err := validation.ValidateStruct(p,
		validation.Field(&p.Email, validation.Required, is.Email),
		validation.Field(&p.Phone, validation.Required, validation.Match(
			regexp.MustCompile(`^[\+]?[(]?[0-9]{3}[)]?[-\s\.]?[0-9]{3}[-\s\.]?[0-9]{4,6}$`))),
		validation.Field(&p.Password, validation.Required, validation.Length(5, 20)),
		validation.Field(&p.Name, validation.Required),
	); err != nil {
		return err
	}

	err := p.HashPassword()
	if err != nil {
		return err
	}
	return nil
}

func (p *Person) HashPassword() error {
	b, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	p.Password = string(b)
	return nil
}
