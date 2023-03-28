package sqlstore

import (
	"database/sql"
	"dip/internal/store"
)

type Store struct {
	db            *sql.DB
	personRep     *PersonRep
	fakeStatusRep *FakeStatusRep
	fakeTaskRep   *FakeTaskRep
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

func (s *Store) Status() store.FakeStatusRepository {
	if s.fakeStatusRep != nil {
		return s.fakeStatusRep
	}

	s.fakeStatusRep = &FakeStatusRep{
		store: s,
	}
	return s.fakeStatusRep
}

func (s *Store) Task() store.FakeTaskRepository {
	if s.fakeTaskRep != nil {
		return s.fakeTaskRep
	}

	s.fakeTaskRep = &FakeTaskRep{
		store: s,
	}
	return s.fakeTaskRep
}
