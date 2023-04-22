package sqlstore

import "dip/internal/models"

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
	for rows.Next() {
		err = rows.Scan(&g.Id, &g.Name, &g.Color)
		if err != nil {
			return nil, err
		}
		groups = append(groups, *g)
	}

	return groups, nil
}
