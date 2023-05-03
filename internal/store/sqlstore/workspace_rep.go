package sqlstore

import (
	"dip/internal/models"
	"time"
)

type WorkspaceRep struct {
	store *Store
}

func (r *WorkspaceRep) GetByUser(id int) (*models.HomePage, error) {
	p := &models.HomePage{
		Ws:       make([]models.WorkspaceJoined, 0),
		Settings: "",
	}
	w := &models.WorkspaceJoined{}

	rows, err := r.store.db.Query(`select pw.workspace_id, w.name, w.description, w.created_at, ur.name, p.settings
		from person_workspace as pw 
		join persons as p on p.id = pw.person_id 
		join user_role as ur on ur.id = pw.role_id 
		join workspaces as w on w.id = pw.workspace_id 
		where p.id = $1 
		order by w.created_at desc`, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&w.Id, &w.Name, &w.Description, &w.CreatedAt, &w.Role, &p.Settings)
		if err != nil {
			return nil, err
		}
		p.Ws = append(p.Ws, *w)
	}

	return p, nil
}

func (r *WorkspaceRep) Create(name, description string) error {
	res, err := r.store.db.Exec(`insert into workspaces
		(name, description, created_at)
		values ($1, $2, $3)`,
		name, description, time.Now(),
	)
	if err != nil {
		return err
	}
	if _, err = res.RowsAffected(); err != nil {
		return err
	}
	return nil
}
func (r *WorkspaceRep) Update(w *models.Workspace) error {
	res, err := r.store.db.Exec(`update workspaces
		set name = $1, description = $2
		where id = $3`,
		w.Name, w.Description,
	)
	if err != nil {
		return err
	}
	if _, err = res.RowsAffected(); err != nil {
		return err
	}
	return nil
}
func (r *WorkspaceRep) Delete(id int) error {
	res, err := r.store.db.Exec(`delete from user_roles ur where ur.id = $1`, id)
	if err != nil {
		return err
	}
	if _, err := res.RowsAffected(); err != nil {
		return err
	}
	return nil
}
