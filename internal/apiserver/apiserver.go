package apiserver

import (
	"database/sql"
	"dip/internal/store/sqlstore"
	"fmt"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Start(config *Config) error {
	db, err := newDB(config.DBconf)
	if err != nil {
		return err
	}
	defer db.Close()
	store := sqlstore.New(db)
	srv := newServer(store)
	srv.logger.Info().Msg("Server started")
	return http.ListenAndServe(config.Port, srv)
}

func newDB(s *DBconfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", fmt.Sprintf("host=%s dbname=%s user=%s password=%s", s.Host, s.DBname, s.User, s.Password))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
