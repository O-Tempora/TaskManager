package sqlstore

import (
	"dip/internal/models"
	"time"
)

type TaskRep struct {
	store *Store
}

func (r *TaskRep) GetAllByGroup(id int) ([]models.TaskOverview, error) {
	t := &models.TaskOverview{}
	var date time.Time
	var ws_id int
	tasks := make([]models.TaskOverview, 0)

	rows, err := r.store.db.Query(`select t.id, t.description, t.created_at, s.name, tg.workspace_id
		from tasks t 
		join statuses s on s.id = t.status_id 
		join task_groups tg on tg.id = t.group_id 
		where t.group_id = $1`, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&t.Id, &t.Description, &date, &t.Status, &ws_id)
		if err != nil {
			return nil, err
		}
		t.CreatedAt = date.Format("2006-01-02")
		t.Executors, err = r.store.Person().GetAllAssignedToTask(t.Id, ws_id)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *t)
	}

	return tasks, nil
}
