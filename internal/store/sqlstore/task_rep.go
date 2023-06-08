package sqlstore

import (
	"dip/internal/models"
	"time"
)

type TaskRep struct {
	store *Store
}

func (r *TaskRep) GetById(taskId int) (*models.Task, error) {
	t := &models.Task{}

	if err := r.store.db.QueryRow(
		`select t.id, t.description, t.created_at, t.start_at, t.finish_at, t.group_id, s.name, t.enddate  from tasks t 
		join statuses s on s.id = t.status_id 
		where t.id = $1`,
		taskId,
	).Scan(&t.Id, &t.Description, &t.CreatedAt, &t.StartdAt, &t.FinishAt, &t.GroupId, &t.Status, &t.EndDate); err != nil {
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
		set description = $1, start_at = $2, finish_at = $3, group_id = $4, status_id = $5, enddate = $6
		where id = $7`,
		task.Description, task.StartdAt, task.FinishAt, task.GroupId, statusId, task.EndDate, task.Id,
	)
	if err != nil {
		return err
	}
	if _, err = res.RowsAffected(); err != nil {
		return err
	}
	return nil
}

func (r *TaskRep) Create(group_id int) (*models.TaskOverview, error) {
	t := &models.Task{
		GroupId:     group_id,
		CreatedAt:   time.Now(),
		StartdAt:    time.Now(),
		FinishAt:    time.Now(),
		Description: "",
		Status:      "In Progress",
	}
	to := &models.TaskOverview{
		Executors:   make([]models.PersonInTask, 0),
		Description: t.Description,
		CreatedAt:   t.CreatedAt.Format("2006-01-02"),
		Status:      t.Status,
	}

	st, err := r.store.Status().GetIdByName(t.Status)
	if err != nil {
		return nil, err
	}

	if err := r.store.db.QueryRow(`insert into tasks
		(description, created_at, start_at, finish_at, group_id, status_id)
		values ($1, $2, $3, $4, $5, $6)
		returning id`,
		t.Description, t.CreatedAt, t.StartdAt, t.FinishAt, t.GroupId, st,
	).Scan(&to.Id); err != nil {
		return nil, err
	}

	return to, nil
}
