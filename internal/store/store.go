package store

type Store interface {
	Person() PersonRepository
	Workspace() WorkspaceRepository
	TaskGroup() TaskGroupRepository
	Status() StatusRepository
	Task() TaskRepository
	Role() RoleRepository
	Comment() CommentRepository
	Invite() InviteRepository
}
