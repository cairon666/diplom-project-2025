package models

type Role struct {
	ID   int32
	Name string
}

type Permission struct {
	ID   int32
	Name string
}

const (
	RoleAdmin             = "admin"
	RoleUser              = "user"
	RoleExternalAppReader = "external_app_reader"
	RoleExternalAppWriter = "external_app_writer"
)
