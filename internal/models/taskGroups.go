package models

type TaskGroup struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type FullGroup struct {
	Group TaskGroup      `json:"group"`
	Tasks []TaskOverview `json:"tasks"`
}
