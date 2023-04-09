package apiserver

import (
	"database/sql"
	"dip/internal/store/sqlstore"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	// connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s", s.Host, s.DBort, s.DBname, s.User, s.Password)
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", s.User, s.Password, s.Host, s.DBort, s.DBname)
	log.Println(connectionString)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		// c, ioErr := ioutil.ReadFile("./migrations/0001init.up.sql")
		// fmt.Println(ioErr)
		// sql := string(c)
		// con := &pgx.Conn{}
		// _, err = con.Exec(context.Background(), sql)
		// log.Println(err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		//log.Println(err)
		return nil, err
	}
	return db, nil
}
