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

func (r *TaskRep) GetById(taskId int) (*models.Task, error) {
	t := &models.Task{}

	if err := r.store.db.QueryRow(
		`select t.id, t.description, t.created_at, t.start_at, t.finish_at, t.group_id, s.name  from tasks t 
		join statuses s on s.id = t.status_id 
		where t.id = $1`,
		taskId,
	).Scan(&t.Id, &t.Description, &t.CreatedAt, &t.StartdAt, &t.FinishAt, &t.GroupId, &t.Status); err != nil {
		return nil, err
	}

	return t, nil
}

func (r *TaskRep) Delete(id int) error {
	res, err := r.store.db.Exec(`delete from tasks t where t.id = $1`, id)
	if err != nil {
		return err
	}

	if _, err := res.RowsAffected(); err != nil {
		return err
	}

	return nil
}

func (r *TaskRep) Update(task *models.Task) error {
	statusId, err := r.store.Status().GetIdByName(task.Status)
	if err != nil {
		return err
	}

	res, err := r.store.db.Exec(`update tasks
		set description = $1, start_at = $2, finish_at = $3, group_id = $4, status_id = $5
		where id = $6`,
		task.Description, task.StartdAt, task.FinishAt, task.GroupId, statusId, task.Id,
	)

	if err != nil {
		return err
	}

	if _, err = res.RowsAffected(); err != nil {
		return err
	}

	return nil
}
