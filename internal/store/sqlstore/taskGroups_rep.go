package sqlstore

import (
	"dip/internal/models"
	"errors"
)

type TaskGroupRep struct {
	store *Store
}

func (r *TaskGroupRep) GetByWorkspaceId(id int) ([]models.TaskGroup, error) {
	g := &models.TaskGroup{}
	groups := make([]models.TaskGroup, 0)

	rows, err := r.store.db.Query(`select tg.Id, tg.Name, tg.Color from task_groups tg where tg.workspace_id = $1 order by tg.Id asc`, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&g.Id, &g.Name, &g.Color)
		if err != nil {
			return nil, err
		}
		groups = append(groups, *g)
	}

	return groups, nil
}

func (r *TaskGroupRep) FindByNameAndWs(ws_id int, name string) (bool, error) {
	var count int
	if err := r.store.db.QueryRow(`select count(*) from task_groups tg 
		where tg.name = $1 and tg.workspace_id = $2`,
		name, ws_id,
	).Scan(&count); err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (r *TaskGroupRep) Create(ws_id int, name, color string) (*models.TG, error) {
	exists, err := r.store.TaskGroup().FindByNameAndWs(ws_id, name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("group with this name already exists")
	}

	tg := &models.TG{}
	if err = r.store.db.QueryRow(`insert into task_groups
		(name, color, workspace_id)
		values ($1, $2, $3)
		returning id, name, color, workspace_id`,
		name, color, ws_id,
	).Scan(&tg.Id, &tg.Name, &tg.Color, &tg.WorkspaceId); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return tg, nil
}

func (r *TaskGroupRep) Delete(id int) error {
	res, err := r.store.db.Exec(`delete from task_groups tg where tg.id = $1`, id)
	if err != nil {
		return err
	}

	if _, err := res.RowsAffected(); err != nil {
		return err
	}

	return nil
}

func (r *TaskGroupRep) Update(tg *models.TG) error {
	res, err := r.store.db.Exec(`update task_groups
		set name = $1, color = $2, workspace_id = $3
		where id = $4`,
		tg.Name, tg.Color, tg.WorkspaceId, tg.Id,
	)
	if err != nil {
		return err
	}
	if _, err = res.RowsAffected(); err != nil {
		return err
	}
	return nil
}
