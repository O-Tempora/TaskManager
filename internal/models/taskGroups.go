package models

type TaskGroup struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type TG struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	WorkspaceId int    `json:"workspace_id"`
}

type FullGroup struct {
	Id    int            `json:"id"`
	Name  string         `json:"name"`
	Color string         `json:"color"`
	Tasks []TaskOverview `json:"tasks"`
}

type GroupPersonal struct {
	Id    int        `json:"group_id"`
	Name  string     `json:"group_name"`
	Tasks []TaskPers `json:"tasks"`
}
