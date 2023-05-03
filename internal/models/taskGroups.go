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
	//Group TaskGroup      `json:"group"`
	Id    int            `json:"id"`
	Name  string         `json:"name"`
	Color string         `json:"color"`
	Tasks []TaskOverview `json:"tasks"`
}
