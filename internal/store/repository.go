package store

import "dip/internal/models"

type PersonRepository interface {
	Create(*models.Person) error
	GetByEmail(email string) (*models.Person, error)
	GetAllAssignedToTask(id int, ws_id int) ([]models.PersonInTask, error)
	GetAllByWorkspace(id int) ([]models.PersonInTask, error)
}

type WorkspaceRepository interface {
	GetByUser(id int) (*models.HomePage, error)
}

type TaskGroupRepository interface {
	GetByWorkspaceId(id int) ([]models.TaskGroup, error)
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
}
