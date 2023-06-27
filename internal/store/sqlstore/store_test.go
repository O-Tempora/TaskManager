package sqlstore

import (
	_ "github.com/lib/pq"
)

var (
	St *Store
)

// func TestMain(m *testing.M) {
// 	cs := "postgres://postgres:postgres@localhost:5432/dip_test?sslmode=disable"
// 	db, err := sql.Open("postgres", cs)
// 	if err != nil {
// 		os.Exit(1)
// 	}
// 	if err = db.Ping(); err != nil {
// 		os.Exit(1)
// 	}

// 	St = New(db)
// 	os.Exit(m.Run())
// }

// func TestStatus_GetAll(t *testing.T) {
// 	res, err := St.Status().GetAll()
// 	if assert.True(t, err == nil, res) {
// 		assert.True(t, len(res) == 3)
// 	}
// }

// func TestStatus_OneByOne(t *testing.T) {
// 	var tests = []struct {
// 		name string
// 		id   int
// 	}{
// 		{"ToDo", 1},
// 		{"Done", 2},
// 		{"Delayed", 3},
// 	}

// 	for _, v := range tests {
// 		id, err := St.Status().GetIdByName(v.name)
// 		if assert.True(t, err == nil) {
// 			assert.Equal(t, v.id, id)
// 		}
// 	}
// }
