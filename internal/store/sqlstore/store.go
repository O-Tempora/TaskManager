package sqlstore

import (
	"database/sql"
	"dip/internal/store"
)

type Store struct {
	db           *sql.DB
	personRep    *PersonRep
	workspaceRep *WorkspaceRep
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
