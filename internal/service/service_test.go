package service

import (
	"database/sql"
	"dip/internal/store/sqlstore"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var (
	Serv *Service
)

func TestMain(m *testing.M) {
	cs := "postgres://postgres:postgres@localhost:5432/dip_test?sslmode=disable"
	db, _ := sql.Open("postgres", cs)
	Serv = New(sqlstore.New(db))
	os.Exit(m.Run())
}
func TestGetAssigned_NoIterations(t *testing.T) {
	res, err := Serv.Logic().GetAllAssignedToTask(13, 1)
	if assert.Nil(t, err) {
		assert.Equal(t, 0, len(res))
	}
}
func TestGetAssigned_TwoAssigned(t *testing.T) {
	res, err := Serv.Logic().GetAllAssignedToTask(33, 1)
	if assert.Nil(t, err) {
		assert.Equal(t, 2, len(res))
	}
}
func TestGetAssigned_OneAssigned(t *testing.T) {
	res, err := Serv.Logic().GetAllAssignedToTask(34, 1)
	if assert.Nil(t, err) {
		assert.Equal(t, 1, len(res))
	}
}
func TestGetAssigned_FailedlQuery(t *testing.T) {
	_, err := Serv.Logic().GetAllAssignedToTask(-1, 1)
	assert.NotNil(t, err)
}
