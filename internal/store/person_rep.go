package store

import "dip/internal/models"

type PersonRep struct {
	store *Store
}

func (r *PersonRep) Create(p *models.Person) (*models.Person, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	if err := r.store.db.QueryRow(
		"INSERT INTO persons (name, password, email, settings, phone) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		p.Name, p.Password, p.Email, p.Settings, p.Phone,
	).Scan(&p.Id); err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PersonRep) GetByEmail(email string) (*models.Person, error) {
	p := &models.Person{}
	if err := r.store.db.QueryRow(
		"SELECT id, name, password, email, settings, phone FROM persons WHERE email = $1",
		email,
	).Scan(&p.Id, &p.Name, &p.Password, &p.Email, &p.Settings, &p.Phone); err != nil {
		return nil, err
	}
	return p, nil
}
