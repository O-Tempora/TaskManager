package store

import "dip/internal/models"

type PersonRepository interface {
	Create(*models.Person) error
	GetByEmail(email string) (*models.Person, error)
	GetIdByName(name string) (int, error)
	GetAllAssignedToTask(id int, ws_id int) ([]models.PersonInTask, error)
	GetAllByWorkspace(id int) ([]models.PersonInTask, error)
	Assign(user string, task int) error
	Dismiss(user string, task int) error
	IsAdmin(user string, ws_id int) (bool, error)
}

type WorkspaceRepository interface {
	GetByUser(id int) (*models.HomePage, error)
	Create(user int, name, description string) (*models.WorkspaceJoined, error)
	Update(w *models.Workspace) error
	Delete(id int) error
	AddUserByEmail(email string, ws_id int) error
}

type TaskGroupRepository interface {
	GetByWorkspaceId(id int) ([]models.TaskGroup, error)
	FindByNameAndWs(ws_id int, name string) (bool, error)
	Create(ws_id int, name, color string) (*models.TG, error)
	Delete(id int) error
	Update(tg *models.TG) error
}

type StatusRepository interface {
	GetAll() ([]models.Status, error)
	GetIdByName(name string) (int, error)
}

type TaskRepository interface {
	GetAllByGroup(id int) ([]models.TaskOverview, error)
	GetById(taskId int) (*models.Task, error)
	Delete(id int) error
	Update(task *models.Task) error
	Create(group_id int) (*models.TaskOverview, error)
}

type RoleRepository interface {
	GetIdByName(name string) (int, error)
	Get(id int) (string, error)
}
