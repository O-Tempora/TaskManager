package store

import "dip/internal/models"

type PersonRepository interface {
	Create(*models.Person) error
	GetByEmail(email string) (*models.Person, error)
}

type WorkspaceRepository interface {
	GetByUser(id int) ([]models.WorkspaceJoined, error)
}
