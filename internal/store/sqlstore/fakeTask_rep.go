package sqlstore

import (
	"dip/internal/models"
	"errors"
	"time"
)

type FakeTaskRep struct {
	store *Store
}

func (r *FakeTaskRep) GetAll() ([]models.JoinedTask, error) {
	buf := models.JoinedTask{}
	var d1, d2, d3 time.Time
	tasks := make([]models.JoinedTask, 0)
	rows, err := r.store.db.Query(`SELECT t.id, t.description, t.created_at, t.start_at, t.finish_at, s.name 
		FROM tasks t JOIN statuses s 
		ON (s.id = t.status_id)`)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&buf.Id, &buf.Description, &d1, &d2, &d3, &buf.Status)
		if err != nil {
			return nil, err
		}
		buf.CreatedAt = d1.Format("2006-01-02")
		buf.StartAt = d2.Format("2006-01-02")
		buf.FinishAt = d3.Format("2006-01-02")
		tasks = append(tasks, buf)
	}

	return tasks, nil
}

func (r *FakeTaskRep) DeleteTask(id int) error {
	_, err := r.store.db.Exec(`DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *FakeTaskRep) Create(task *models.JoinedTask, status_id int) error {
	date1, err := time.Parse("2006-01-02", task.StartAt)
	if err != nil {
		return err
	}
	date2, err := time.Parse("2006-01-02", task.FinishAt)
	if err != nil {
		return err
	}

	if date1.After(date2) {
		return errors.New("Start date can not be earlier than finish date")
	}

	date3 := time.Now().Format("2006-01-02")

	_, err = r.store.db.Exec(`INSERT INTO tasks 
		(description, start_at, finish_at, created_at, group_id, status_id)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		task.Description, date1, date2, date3, 1, status_id)

	return err
}

func (r *FakeTaskRep) Update(task *models.JoinedTask, status_id int) error {
	date1, err := time.Parse("2006-01-02", task.StartAt)
	if err != nil {
		return err
	}
	date2, err := time.Parse("2006-01-02", task.FinishAt)
	if err != nil {
		return err
	}

	if date1.After(date2) {
		return errors.New("Start date can not be earlier than finish date")
	}

	_, err = r.store.db.Exec(`UPDATE tasks 
		SET description=$1,
			start_at=$2,
			finish_at=$3,
			status_id=$4
		WHERE id = $5`,
		task.Description, date1, date2, status_id, task.Id)

	return err
}

func (r *FakeTaskRep) Get(id int) (models.JoinedTask, error) {
	buf := models.JoinedTask{}
	var d1, d2, d3 time.Time
	rows, err := r.store.db.Query(`SELECT t.id, t.description, t.created_at, t.start_at, t.finish_at, s.name 
		FROM tasks t JOIN statuses s 
		ON (s.id = t.status_id)
		WHERE t.id=$1`, id)

	if err != nil {
		return buf, err
	}

	for rows.Next() {
		err = rows.Scan(&buf.Id, &buf.Description, &d1, &d2, &d3, &buf.Status)
		if err != nil {
			return buf, err
		}
		buf.CreatedAt = d1.Format("2006-01-02")
		buf.StartAt = d2.Format("2006-01-02")
		buf.FinishAt = d3.Format("2006-01-02")
	}

	return buf, nil
}
