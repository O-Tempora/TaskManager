package sqlstore

import "errors"

type RoleRep struct {
	store *Store
}

func (r *RoleRep) GetIdByName(name string) (int, error) {
	id := -1
	if err := r.store.db.QueryRow(`select ur.id from user_role ur
		where ur.name = $1`,
		name,
	).Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

func (r *RoleRep) Get(id int) (string, error) {
	name := ""
	if err := r.store.db.QueryRow(`select ur.name from user_role ur
		where ur.id = $1`,
		id,
	).Scan(&name); err != nil {
		return "", err
	}
	if name == "" {
		return "", errors.New("Role doesn't exist")
	}
	return name, nil
}
