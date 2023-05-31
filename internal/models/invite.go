package models

import "time"

type Invite struct {
	Id         int       `json:"id"`
	SenderId   int       `json:"sender"`
	ReceiverId int       `json:"receiver"`
	WsId       int       `json:"ws"`
	SentAt     time.Time `json:"sent_at"`
}

type InviteShow struct {
	Id          int       `json:"id"`
	Sender      string    `json:"sender"`
	ReceiverId  int       `json:"receiver"`
	SentAt      time.Time `json:"sent_at"`
	Workspace   string    `json:"workspace"`
	WorkspaceId int       `json:"ws_id"`
}
