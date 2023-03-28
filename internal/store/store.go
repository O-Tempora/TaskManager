package store

type Store interface {
	Person() PersonRepository
	Status() FakeStatusRepository
	Task() FakeTaskRepository
}
