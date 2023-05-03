package sqlstore

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
