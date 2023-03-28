package store

import "dip/internal/models"

type PersonRepository interface {
	Create(*models.Person) error
	GetByEmail(email string) (*models.Person, error)
}

type FakeStatusRepository interface {
	GetAll() ([]models.Status, error)
	GetIdByName(name string) (int, error)
}

type FakeTaskRepository interface {
	GetAll() ([]models.JoinedTask, error)
}
