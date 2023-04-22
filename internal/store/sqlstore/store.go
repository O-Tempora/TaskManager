package sqlstore

import (
	"database/sql"
	"dip/internal/store"
)

type Store struct {
	db           *sql.DB
	personRep    *PersonRep
	workspaceRep *WorkspaceRep
	taskGroupRep *TaskGroupRep
	statusRep    *StatusRep
	taskRep      *TaskRep
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Person() store.PersonRepository {
	if s.personRep != nil {
		return s.personRep
	}

	s.personRep = &PersonRep{
		store: s,
	}
	return s.personRep
}

func (s *Store) Workspace() store.WorkspaceRepository {
	if s.workspaceRep != nil {
		return s.workspaceRep
	}

	s.workspaceRep = &WorkspaceRep{
		store: s,
	}
	return s.workspaceRep
}

func (s *Store) TaskGroup() store.TaskGroupRepository {
	if s.taskGroupRep != nil {
		return s.taskGroupRep
	}

	s.taskGroupRep = &TaskGroupRep{
		store: s,
	}
	return s.taskGroupRep
}

func (s *Store) Status() store.StatusRepository {
	if s.statusRep != nil {
		return s.statusRep
	}

	s.statusRep = &StatusRep{
		store: s,
	}
	return s.statusRep
}

func (s *Store) Task() store.TaskRepository {
	if s.taskRep != nil {
		return s.taskRep
	}

	s.taskRep = &TaskRep{
		store: s,
	}
	return s.taskRep
}
