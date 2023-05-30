package models

import "time"

type Comment struct {
	Id        int       `json:"id"`
	Content   string    `json:"content"`
	PersonId  int       `json:"person_id"`
	TaskId    int       `json:"task_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentShow struct {
	Id        int       `json:"id"`
	Content   string    `json:"content"`
	Person    string    `json:"person"`
	PersonId  int       `json:"person_id"`
	CreatedAt time.Time `json:"created_at"`
}
