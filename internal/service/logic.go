package service

import (
	"dip/internal/middleware"
	"dip/internal/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Logic interface {
	SendInvite(email string, ws_id, user_id int) error
	DeclineInvite(invite_id int) error
	AcceptInvite(invite_id, ws_id, usr_id int) (*models.WorkspaceJoined, error)

	AddToWsByEmail(email string, ws_id int) error
	GetWsByUser(id int) (*models.HomePage, error)
	GetHome(id int) (*models.HomePage, error)
	GetFullWorkspace(ws_id int) (*models.WorkspaceFull, error)

	LogIn(email, password string) (string, error)
	SignUp(name, email, phone, password string) error
	GetAllAssignedToTask(id int, ws_id int) ([]models.PersonInTask, error)
	Assign(name string, task int) error
	Dismiss(name string, task int) error
	IsAdmin(name string, ws_id int) (bool, error)
	LeaveWs(id, ws_id, next_admin_id int) error

	GetAllTasksByGroup(id int) ([]models.TaskOverview, error)
	MoveTask(id, gr_id int) error
	GetAllTasksByUser(id int) ([]models.PersonalTasksInWs, error)
}

type LogicUnit struct {
	service *Service
}

func (l *LogicUnit) SendInvite(email string, ws_id, user_id int) error {
	p, err := l.service.store.Person().GetByEmail(email)
	if err != nil {
		return err
	}
	err = l.service.store.Invite().Create(&models.Invite{
		Id:         p.Id,
		SenderId:   user_id,
		ReceiverId: p.Id,
		WsId:       ws_id,
		SentAt:     time.Now(),
	})
	return err
}
func (l *LogicUnit) DeclineInvite(invite_id int) error {
	if err := l.service.store.Invite().Delete(invite_id); err != nil {
		return err
	}
	return nil
}
func (l *LogicUnit) AcceptInvite(invite_id, ws_id, usr_id int) (*models.WorkspaceJoined, error) {
	role_id, err := l.service.store.Role().GetIdByName("User")
	if err != nil {
		return nil, err
	}
	res, err := l.service.store.DB().Exec(`insert into person_workspace
		(person_id, workspace_id, role_id)
		values($1, $2, $3) on conflict do nothing`,
		usr_id, ws_id, role_id,
	)
	if err != nil {
		return nil, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if err := l.service.store.Invite().Delete(invite_id); err != nil {
		return nil, err
	}

	w := &models.WorkspaceJoined{}
	if err := l.service.store.DB().QueryRow(`select pw.workspace_id, w.name, w.description, w.created_at, w.isactive, w.closed_at, ur.name
		from person_workspace as pw 
		join persons as p on p.id = pw.person_id 
		join user_role as ur on ur.id = pw.role_id 
		join workspaces as w on w.id = pw.workspace_id 
		where w.id = $1 `,
		ws_id).Scan(&w.Id, &w.Name, &w.Description, &w.CreatedAt, &w.IsActive, &w.ClosedAt, &w.Role); err != nil {
		return nil, err
	}

	return w, nil
}

func (l *LogicUnit) AddToWsByEmail(email string, ws_id int) error {
	p, err := l.service.store.Person().GetByEmail(email)
	if err != nil {
		return err
	}
	role_id, err := l.service.store.Role().GetIdByName("User")
	if err != nil {
		return err
	}
	res, err := l.service.store.DB().Exec(`insert into person_workspace
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
func (l *LogicUnit) GetWsByUser(id int) (*models.HomePage, error) {
	p := &models.HomePage{
		Ws:       make([]models.WorkspaceJoined, 0),
		Settings: "",
	}
	w := &models.WorkspaceJoined{}

	rows, err := l.service.store.DB().Query(`select pw.workspace_id, w.name, w.description, w.created_at, w.isactive, w.closed_at, ur.name, p.settings
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
func (l *LogicUnit) GetHome(id int) (*models.HomePage, error) {
	ws, err := l.GetWsByUser(id)
	if err != nil {
		return nil, err
	}

	return ws, nil
}
func (l *LogicUnit) GetFullWorkspace(ws_id int) (*models.WorkspaceFull, error) {
	ws, err := l.service.store.Workspace().GetById(ws_id)
	if err != nil {
		return nil, err
	}
	fg := &models.FullGroup{}
	groups := make([]models.FullGroup, 0)

	tgs, err := l.service.store.TaskGroup().GetByWorkspaceId(ws_id)
	if err != nil {
		return nil, err
	}

	for _, v := range tgs {
		fg.Id = v.Id
		fg.Color = v.Color
		fg.Name = v.Name
		fg.Tasks, err = l.GetAllTasksByGroup(v.Id)
		if err != nil {
			return nil, err
		}
		groups = append(groups, *fg)
	}

	wsf := &models.WorkspaceFull{
		WS:     *ws,
		Groups: groups,
	}
	return wsf, nil
}

func (l *LogicUnit) LogIn(email, password string) (string, error) {
	p, err := l.service.store.Person().GetByEmail(email)
	if err != nil {
		return "", err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(password)); err != nil {
		return "", err
	}

	token, err := middleware.GenerateJWT(p)
	if err != nil {
		return "", err
	}

	return token, nil
}
func (l *LogicUnit) SignUp(name, email, phone, password string) error {
	p := &models.Person{
		Email:        email,
		Password:     password,
		Name:         name,
		Phone:        phone,
		Settings:     "",
		IsMaintainer: false,
	}

	if err := l.service.store.Person().Create(p); err != nil {
		return err
	}

	return nil
}

func (l *LogicUnit) GetAllAssignedToTask(id int, ws_id int) ([]models.PersonInTask, error) {
	p := &models.PersonInTask{}
	persons := make([]models.PersonInTask, 0)

	//Mb add n1.name to returned values of select
	rows, err := l.service.store.DB().Query(`select p.id, p.name, n1.name as role from persons p 
		join (select * from person_workspace pw 
			join user_role ur
			on pw.role_id = ur.id 
			where pw.workspace_id = $1
		) as n1 on p.id = n1.person_id 
		where p.id in (
			select pt.person_id from person_task pt 
			where pt.task_id = $2
		)`, ws_id, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&p.Id, &p.Name, &p.Role)
		if err != nil {
			return nil, err
		}
		persons = append(persons, *p)
	}

	return persons, nil
}
func (l *LogicUnit) Assign(name string, task int) error {
	id, err := l.service.store.Person().GetIdByName(name)
	if err != nil {
		return err
	}
	res, err := l.service.store.DB().Exec(`insert into person_task 
		(person_id, task_id)
		values ($1, $2)
		on conflict do nothing`,
		id, task,
	)
	if err != nil {
		return err
	}
	if _, err = res.RowsAffected(); err != nil {
		return err
	}
	return nil
}
func (l *LogicUnit) Dismiss(name string, task int) error {
	id, err := l.service.store.Person().GetIdByName(name)
	if err != nil {
		return err
	}
	res, err := l.service.store.DB().Exec(`delete from person_task pt
		where pt.person_id = $1 and pt.task_id = $2`,
		id, task,
	)
	if err != nil {
		return err
	}
	if _, err = res.RowsAffected(); err != nil {
		return err
	}
	return nil
}
func (l *LogicUnit) IsAdmin(name string, ws_id int) (bool, error) {
	id, err := l.service.store.Person().GetIdByName(name)
	if err != nil {
		return false, err
	}
	role_id := -1
	if err = l.service.store.DB().QueryRow(`select pw.role_id from person_workspace pw
		where pw.person_id = $1 and pw.workspace_id = $2`,
		id, ws_id,
	).Scan(&role_id); err != nil {
		return false, err
	}

	role, err := l.service.store.Role().Get(role_id)
	if err != nil {
		return false, err
	}
	if role == "Admin" {
		return true, nil
	}
	return false, nil
}
func (l *LogicUnit) LeaveWs(id, ws_id, next_admin_id int) error {
	res, err := l.service.store.DB().Exec(`delete from person_workspace where person_id = $1 and workspace_id = $2`, id, ws_id)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	//no next admin provided
	if next_admin_id == -1 {
		return nil
	}

	roleId, err := l.service.store.Role().GetIdByName("Admin")
	if err != nil {
		return err
	}
	res, err = l.service.store.DB().Exec(`update person_workspace set role_id = $1 where person_id = $2 and workspace_id = $3`, roleId, next_admin_id, ws_id)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	return err
}

func (l *LogicUnit) GetAllTasksByGroup(id int) ([]models.TaskOverview, error) {
	t := &models.TaskOverview{}
	var date time.Time
	var ws_id int
	tasks := make([]models.TaskOverview, 0)

	rows, err := l.service.store.DB().Query(`select t.id, t.description, t.created_at, s.name, tg.workspace_id, t.enddate
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
		err = rows.Scan(&t.Id, &t.Description, &date, &t.Status, &ws_id, &t.EndDate)
		if err != nil {
			return nil, err
		}
		t.CreatedAt = date.Format("2006-01-02")
		t.Executors, err = l.GetAllAssignedToTask(t.Id, ws_id)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *t)
	}

	return tasks, nil
}
func (l *LogicUnit) MoveTask(id, gr_id int) error {
	res, err := l.service.store.DB().Exec(`update tasks set group_id = $1 where id = $2`, gr_id, id)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	return err
}
func (l *LogicUnit) GetAllTasksByUser(id int) ([]models.PersonalTasksInWs, error) {
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
		EndDate     *time.Time
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

	//q.status_id <> 2 mean that task is not done yet
	res := make([]models.PersonalTasksInWs, 0)
	rows, err := l.service.store.DB().Query(`select w.id, w."name", tg.id, tg."name", q.id, q.description, q.created_at, q.start_at, q.finish_at, q."name", q.enddate from workspaces w
		join task_groups tg on tg.workspace_id = w.id 
		join (select t.id, t.description, t.created_at, t.start_at, t.finish_at, t.group_id, t.status_id, t.enddate, s."name" from tasks t 
			join statuses s on s.id = t.status_id) as q on q.group_id = tg.id 
		where (q.status_id <> 2
		) and w.id in (
			select pw.workspace_id from person_workspace pw where pw.person_id = $1
		) and q.id in (
			select pt.task_id from person_task pt where pt.person_id = $2
		) order by w.id`, id, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&tp.Ws_id, &tp.Ws_name, &tp.Group_id, &tp.Group_name, &tp.Task_id, &tp.Description, &tp.CreatedAt, &tp.StartAt, &tp.FinishAt, &tp.Status, &tp.EndDate)
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
			EndDate:     tp.EndDate,
			Status:      tp.Status,
		})

	}
	if rows.Err() != nil {
		return nil, err
	}

	return res, nil
}
