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
	srv.logger.Info().Msgf("Server started at port %s", config.Port)
	return http.ListenAndServe(config.Port, srv)
}

func newDB(s *DBconfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", fmt.Sprintf("host=%s dbname=%s user=%s password=%s", s.Host, s.DBname, s.User, s.Password))
	if err != nil {
		// c, ioErr := ioutil.ReadFile("./migrations/dump-dip-202304081654.sql")
		// fmt.Println(ioErr)
		// sql := string(c)
		// con := &pgx.Conn{}
		// con.Exec(context.Background(), sql)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
