package models

import (
	"time"
)

type Task struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	StartdAt    time.Time `json:"startAt"`
	FinishAt    time.Time `json:"finishAt"`
	GroupId     int       `json:"groupId"`
	Status      string    `json:"statusId"`
}

type TaskOverview struct {
	Id          int            `json:"id"`
	Description string         `json:"description"`
	CreatedAt   string         `json:"createdAt"`
	Status      string         `json:"status"`
	Executors   []PersonInTask `json:"persons"`
}
