package memory

import "time"

type Role string

const (
	Member Role = "Member"
	Admin  Role = "Admin"
)

type User struct {
	ID         string
	Name       string
	SecondName string
	Role       Role
	Created    time.Time
}
