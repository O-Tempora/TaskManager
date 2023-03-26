package store

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	//_ "github.com/lib/pq"
)

type Store struct {
	config    *Config
	db        *sql.DB
	personRep *PersonRep
}

func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

func (s *Store) Open() error {
	db, err := sql.Open("pgx", fmt.Sprintf("host=%s dbname=%s user=%s password=%s", s.config.Host, s.config.DBname, s.config.User, s.config.Password))
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) Person() *PersonRep {
	if s.personRep != nil {
		return s.personRep
	}

	s.personRep = &PersonRep{
		store: s,
	}
	return s.personRep
}
