package sqlstore

import (
	"dip/internal/models"
	"errors"
)

type PersonRep struct {
	store *Store
}

func (r *PersonRep) Create(p *models.Person) error {
	if err := p.Validate(); err != nil {
		return err
	}
	return r.store.db.QueryRow(
		"INSERT INTO persons (name, password, email, settings, phone, ismaintainer) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		p.Name, p.Password, p.Email, p.Settings, p.Phone, p.IsMaintainer,
	).Scan(&p.Id)
}

func (r *PersonRep) GetByEmail(email string) (*models.Person, error) {
	p := &models.Person{}
	if err := r.store.db.QueryRow(
		"SELECT id, name, password, email, settings, phone, ismaintainer FROM persons WHERE email = $1",
		email,
	).Scan(&p.Id, &p.Name, &p.Password, &p.Email, &p.Settings, &p.Phone, &p.IsMaintainer); err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PersonRep) GetAllByWorkspace(id int) ([]models.PersonWS, error) {
	p := &models.PersonWS{}
	persons := make([]models.PersonWS, 0)

	rows, err := r.store.db.Query(`select p.id, p.name, n1.name as role, p.email, p.phone from persons p 
		join (select * from person_workspace pw 
			join user_role ur
			on pw.role_id = ur.id 
			where pw.workspace_id = $1
		) as n1 on p.id = n1.person_id `, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&p.Id, &p.Name, &p.Role, &p.Email, &p.Phone)
		if err != nil {
			return nil, err
		}
		persons = append(persons, *p)
	}

	return persons, nil
}

func (r *PersonRep) GetIdByName(name string) (int, error) {
	id := -1
	if err := r.store.db.QueryRow(`select p.id from persons p
		where p.name = $1`,
		name,
	).Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

func (r *PersonRep) Delete(id int) error {
	rows, err := r.store.db.Exec(`delete from persons where id = $1 and ismaintainer = false`, id)
	if err != nil {
		return err
	}
	ra, err := rows.RowsAffected()
	if err != nil {
		return err
	}
	if ra == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (r *PersonRep) Update(id int, p models.Person) error {
	rows, err := r.store.db.Exec(`update persons 
		set name = $1, email = $2, settings = $3, phone = $4
		where id = $5`,
		p.Name, p.Email, p.Settings, p.Phone, id,
	)
	if err != nil {
		return err
	}
	ra, err := rows.RowsAffected()
	if err != nil {
		return err
	}
	if ra == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (r *PersonRep) GetAll(page, take int) ([]models.PersonShow, error) {
	res := make([]models.PersonShow, 0)
	ps := models.PersonShow{}

	rows, err := r.store.db.Query(`select p.id, p.name, p.email, p.phone, p.ismaintainer from persons p limit $1 offset $2`, take, (page-1)*take)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&ps.Id, &ps.Name, &ps.Email, &ps.Phone, &ps.IsMaintainer)
		if err != nil {
			return nil, err
		}
		res = append(res, ps)
	}
	return res, nil
}

func (r *PersonRep) GetNameById(id int) (string, error) {
	var name string
	if err := r.store.db.QueryRow(`select p.name from persons p where p.id = $1`, id).Scan(&name); err != nil {
		return "", err
	}
	return name, nil
}
