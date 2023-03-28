package sqlstore

import "dip/internal/models"

type FakeTaskRep struct {
	store *Store
}

func (r *FakeTaskRep) GetAll() ([]models.JoinedTask, error) {
	buf := models.JoinedTask{}
	tasks := make([]models.JoinedTask, 0)
	rows, err := r.store.db.Query(`SELECT t.id, t.description, t.created_at, t.start_at, t.finish_at, s.name 
		FROM tasks t JOIN statuses s 
		ON (s.id = t.status_id)`)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&buf.Id, &buf.Description, &buf.CreatedAt, &buf.StartAt, &buf.FinishAt, &buf.Status)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, buf)
	}

	return tasks, nil
}
