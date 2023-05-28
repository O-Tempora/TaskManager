package models

import "time"

type Workspace struct {
	Id          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	IsActive    bool       `json:"isActive"`
	ClosedAt    *time.Time `json:"closed_at"`
}

type WorkspaceJoined struct {
	Id          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	IsActive    bool       `json:"isActive"`
	ClosedAt    *time.Time `json:"closed_at"`
	Role        string     `json:"role"`
}

type HomePage struct {
	Ws       []WorkspaceJoined `json:"workspaces"`
	Settings string            `json:"settings"`
}

type WorkspaceFull struct {
	WS     Workspace   `json:"info"`
	Groups []FullGroup `json:"groups"`
}
