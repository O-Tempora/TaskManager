package sqlstore

import (
	"dip/internal/models"
)

type PersonRep struct {
	store *Store
}

func (r *PersonRep) Create(p *models.Person) error {
	if err := p.Validate(); err != nil {
		return err
	}
	return r.store.db.QueryRow(
		"INSERT INTO persons (name, password, email, settings, phone, ismaintainer) VALUES ($1, $2, $3, $4, $5) RETURNING id",
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

func (r *PersonRep) GetAllAssignedToTask(id int, ws_id int) ([]models.PersonInTask, error) {
	p := &models.PersonInTask{}
	persons := make([]models.PersonInTask, 0)

	//Mb add n1.name to returned values of select
	rows, err := r.store.db.Query(`select p.id, p.name as role from persons p 
		join (select * from person_workspace pw 
			join user_role ur
			on pw.role_id = ur.id 
			where pw.workspace_id = $1
		) as n1 on p.id = n1.person_id 
		where p.id in (
			select pt.person_id from person_task pt 
			where pt.task_id = $2
		)`, id, ws_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&p.Id, &p.Name)
		if err != nil {
			return nil, err
		}
		persons = append(persons, *p)
	}

	return persons, nil
}
