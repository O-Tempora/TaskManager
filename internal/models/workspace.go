package models

import "time"

type Workspace struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type WorkspaceJoined struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	Role        string    `json:"role"`
}

type HomePage struct {
	Ws       []WorkspaceJoined `json:"workspaces"`
	Settings string            `json:"settings"`
}

type WorkspaceFull struct {
	Id     int         `json:"id"`
	Groups []FullGroup `json:"groups"`
}
