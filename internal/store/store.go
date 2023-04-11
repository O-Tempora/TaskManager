package store

type Store interface {
	Person() PersonRepository
	Workspace() WorkspaceRepository
}
