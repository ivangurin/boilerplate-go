package model

type Action string

const (
	ActionCreate Action = "create"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
	ActionLogin  Action = "login"
	ActionLogout Action = "logout"
)
