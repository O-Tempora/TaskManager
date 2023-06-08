package store

import "dip/internal/models"

type PersonRepository interface {
	Create(*models.Person) error
	GetByEmail(email string) (*models.Person, error)
	GetIdByName(name string) (int, error)
	GetAllByWorkspace(id int) ([]models.PersonInTask, error)
	GetNameById(id int) (string, error)
	GetAll(page, take int) ([]models.PersonShow, error)
	Delete(id int) error
	Update(id int, p models.Person) error
}

type WorkspaceRepository interface {
	Create(user int, name, description string) (*models.WorkspaceJoined, error)
	Update(w *models.Workspace, id int) error
	Delete(id int) error
	GetById(id int) (*models.Workspace, error)
	GetAll(page, take int) ([]models.Workspace, error)
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
	GetById(taskId int) (*models.Task, error)
	Delete(id int) error
	Update(task *models.Task) error
	Create(group_id int) (*models.TaskOverview, error)
}

type RoleRepository interface {
	GetIdByName(name string) (int, error)
	Get(id int) (string, error)
}

type CommentRepository interface {
	GetByTask(id int) ([]models.CommentShow, error)
	Create(c models.Comment) (*models.CommentShow, error)
	Delete(id int) error
}

type InviteRepository interface {
	GetAll(user_id int) ([]models.InviteShow, error)
	Create(inv *models.Invite) error
	Delete(invite_id int) error
}
