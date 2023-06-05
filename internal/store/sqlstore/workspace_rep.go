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

	rows, err := r.store.db.Query(`select pw.workspace_id, w.name, w.description, w.created_at, w.isactive, w.closed_at, ur.name, p.settings
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
		err = rows.Scan(&w.Id, &w.Name, &w.Description, &w.CreatedAt, &w.IsActive, &w.ClosedAt, &w.Role, &p.Settings)
		if err != nil {
			return nil, err
		}
		p.Ws = append(p.Ws, *w)
	}

	return p, nil
}

func (r *WorkspaceRep) Create(user int, name, description string) (*models.WorkspaceJoined, error) {
	var id int = -1
	if err := r.store.db.QueryRow(`insert into workspaces
		(name, description, created_at, isactive)
		values ($1, $2, $3, $4)
		returning id`,
		name, description, time.Now(), true,
	).Scan(&id); err != nil {
		return nil, err
	}

	role_id, err := r.store.Role().GetIdByName("Admin")
	if err != nil {
		return nil, err
	}

	ws := &models.WorkspaceJoined{
		Id:          id,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		Role:        "Admin",
		IsActive:    true,
		ClosedAt:    nil,
	}

	_, err = r.store.db.Exec(`insert into person_workspace
		(person_id, workspace_id, role_id)
		values ($1, $2, $3)
		on conflict do nothing`,
		user, id, role_id,
	)
	if err != nil {
		return nil, err
	}

	return ws, nil
}
func (r *WorkspaceRep) Update(w *models.Workspace, id int) error {
	res, err := r.store.db.Exec(`update workspaces
		set name = $1, description = $2, isactive = $3, closed_at = $4
		where id = $5`,
		w.Name, w.Description, w.IsActive, w.ClosedAt, id,
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
	res, err := r.store.db.Exec(`delete from workspaces where id = $1`, id)
	if err != nil {
		return err
	}
	if _, err := res.RowsAffected(); err != nil {
		return err
	}
	return nil
}

func (r *WorkspaceRep) AddUserByEmail(email string, ws_id int) error {
	p, err := r.store.personRep.GetByEmail(email)
	if err != nil {
		return err
	}
	role_id, err := r.store.Role().GetIdByName("User")
	if err != nil {
		return err
	}
	res, err := r.store.db.Exec(`insert into person_workspace
		(person_id, workspace_id, role_id)
		values($1, $2, $3) on conflict do nothing`,
		p.Id, ws_id, role_id,
	)
	if err != nil {
		return err
	}
	if _, err = res.RowsAffected(); err != nil {
		return err
	}
	return nil
}

func (r *WorkspaceRep) GetById(id int) (*models.Workspace, error) {
	ws := &models.Workspace{}
	if err := r.store.db.QueryRow(`select * from workspaces ws where ws.id = $1`, id).
		Scan(&ws.Id, &ws.Name, &ws.Description, &ws.CreatedAt, &ws.IsActive, &ws.ClosedAt); err != nil {
		return nil, err
	}
	return ws, nil
}

func (r *WorkspaceRep) GetAll(page, take int) ([]models.Workspace, error) {
	rows, err := r.store.db.Query(`select * from workspaces limit $1 offset $2`, take, (page-1)*take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]models.Workspace, 0)
	ws := models.Workspace{}
	for rows.Next() {
		err = rows.Scan(&ws.Id, &ws.Name, &ws.Description, &ws.CreatedAt, &ws.IsActive, &ws.ClosedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, ws)
	}
	return res, rows.Err()
}
