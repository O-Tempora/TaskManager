package models

import (
	"time"
)

type Task struct {
	Id          int        `json:"id"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	StartdAt    time.Time  `json:"startAt"`
	FinishAt    time.Time  `json:"finishAt"`
	EndDate     *time.Time `json:"endDate"`
	GroupId     int        `json:"groupId"`
	Status      string     `json:"status"`
}

type TaskOverview struct {
	Id          int            `json:"id"`
	Description string         `json:"description"`
	FinishAt    string         `json:"finishAt"`
	Status      string         `json:"status"`
	EndDate     *time.Time     `json:"endDate"`
	Executors   []PersonInTask `json:"persons"`
}
type TaskPers struct {
	Id          int        `json:"id"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	StartAt     time.Time  `json:"startAt"`
	FinishAt    time.Time  `json:"finishAt"`
	EndDate     *time.Time `json:"endDate"`
	Status      string     `json:"status"`
}

type PersonalTasksInWs struct {
	Id     int             `json:"ws_id"`
	Name   string          `json:"ws_name"`
	Groups []GroupPersonal `json:"groups"`
}
