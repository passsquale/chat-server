package model

type AccessCheck struct {
	Roles []UserRole
}

type UserRole int8

const (
	UNKNOWN UserRole = iota
	USER
	ADMIN
)
