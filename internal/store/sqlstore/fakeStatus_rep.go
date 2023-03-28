package sqlstore

import "dip/internal/models"

type FakeStatusRep struct {
	store *Store
}

func (r *FakeStatusRep) GetAll() ([]models.Status, error) {
	buf := models.Status{}
	st := make([]models.Status, 0)
	rows, err := r.store.db.Query("SELECT id, name FROM statuses")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&buf.Id, &buf.Name)
		if err != nil {
			return nil, err
		}
		st = append(st, buf)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return st, nil
}

func (r *FakeStatusRep) GetIdByName(name string) (int, error) {
	var buf int
	rows, err := r.store.db.Query("SELECT id FROM statuses WHERE name=$1 LIMIT 1", name)
	if err != nil {
		return -1, err
	}

	for rows.Next() {
		err = rows.Scan(&buf)
		if err != nil {
			return -1, err
		}
	}

	return buf, nil
}
