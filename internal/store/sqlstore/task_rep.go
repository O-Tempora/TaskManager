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
		where t.group_id = $1
		order by t.id desc`, id)

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

func (r *TaskRep) GetAllByUser(id int) ([]models.PersonalTasksInWs, error) {
	tp := struct {
		Ws_id       int
		Ws_name     string
		Group_id    int
		Group_name  string
		Task_id     int
		Description string
		CreatedAt   time.Time
		StartAt     time.Time
		FinishAt    time.Time
		Status      string
	}{}

	//Checks if workspace with these ID and NAME was already added to slice. Returns index
	wsExists := func(s []models.PersonalTasksInWs, id int, name string) (bool, int) {
		for i, v := range s {
			if v.Id == id && v.Name == name {
				return true, i
			}
		}
		return false, -1
	}
	//Checks if group with these ID and NAME was already added to slice. Returns index
	groupExists := func(s []models.GroupPersonal, id int, name string) (bool, int) {
		for i, v := range s {
			if v.Id == id && v.Name == name {
				return true, i
			}
		}
		return false, -1
	}

	res := make([]models.PersonalTasksInWs, 0)
	rows, err := r.store.db.Query(`select w.id, w."name", tg.id, tg."name", q.id, q.description, q.created_at, q.start_at, q.finish_at, q."name" from workspaces w
		join task_groups tg on tg.workspace_id = w.id 
		join (select t.id, t.description, t.created_at, t.start_at, t.finish_at, t.group_id, s."name" from tasks t 
			join statuses s on s.id = t.status_id) as q on q.group_id = tg.id 
		where w.id in (
			select pw.workspace_id from person_workspace pw where pw.person_id = $1
		) and q.id in (
			select pt.task_id from person_task pt where pt.person_id = $2
		) order by w.id`, id, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&tp.Ws_id, &tp.Ws_name, &tp.Group_id, &tp.Group_name, &tp.Task_id, &tp.Description, &tp.CreatedAt, &tp.StartAt, &tp.FinishAt, &tp.Status)
		if err != nil {
			return nil, err
		}

		//Add new ws to slice if needed
		ok, ind := wsExists(res, tp.Ws_id, tp.Ws_name)
		if !ok {
			pt := &models.PersonalTasksInWs{
				Id:     tp.Ws_id,
				Name:   tp.Ws_name,
				Groups: make([]models.GroupPersonal, 0),
			}
			res = append(res, *pt)
			ind = len(res) - 1
		}

		ok, grInd := groupExists(res[ind].Groups, tp.Group_id, tp.Group_name)
		if !ok {
			gr := &models.GroupPersonal{
				Id:    tp.Group_id,
				Name:  tp.Group_name,
				Tasks: make([]models.TaskPers, 0),
			}
			res[ind].Groups = append(res[ind].Groups, *gr)
			grInd = len(res[ind].Groups) - 1
		}

		res[ind].Groups[grInd].Tasks = append(res[ind].Groups[grInd].Tasks, models.TaskPers{
			Id:          tp.Task_id,
			Description: tp.Description,
			CreatedAt:   tp.CreatedAt,
			StartAt:     tp.StartAt,
			FinishAt:    tp.FinishAt,
			Status:      tp.Status,
		})

	}
	if rows.Err() != nil {
		return nil, err
	}

	return res, nil
}
