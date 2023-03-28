package models

import "time"

type Task struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	StartdAt    time.Time `json:"startAt"`
	FinishAt    time.Time `json:"finishAt"`
	GroupId     int       `json:"groupId"`
	StatusId    int       `json:"statusId"`
}

type JoinedTask struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	StartAt     time.Time `json:"startAt"`
	FinishAt    time.Time `json:"finishAt"`
	Status      string    `json:"status"`
}
