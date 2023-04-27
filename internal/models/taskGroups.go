package models

type TaskGroup struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type FullGroup struct {
	//Group TaskGroup      `json:"group"`
	Id    int            `json:"id"`
	Name  string         `json:"name"`
	Color string         `json:"color"`
	Tasks []TaskOverview `json:"tasks"`
}
