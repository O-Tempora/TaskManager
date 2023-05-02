package sqlstore

import "dip/internal/models"

type StatusRep struct {
	store *Store
}

func (r *StatusRep) GetAll() ([]models.Status, error) {
	s := &models.Status{}
	statuses := make([]models.Status, 0)

	rows, err := r.store.db.Query(`select * from statuses`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&s.Id, &s.Name)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, *s)
	}

	return statuses, nil
}

func (r *StatusRep) GetIdByName(name string) (int, error) {
	var id int = 0

	if err := r.store.db.QueryRow(`select s.Id from statuses s where s.Name = $1`, name).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
