package store

type Store interface {
	Person() PersonRepository
}
