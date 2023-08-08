package roles

import "strings"

type Role int

const (
	Trainee Role = iota
	Admin
)

func GetRoleFromString(s string) Role {
	switch strings.ToLower(s) {
	case "trainee":
		return Trainee
	case "admin":
		return Admin
	default:
		return -1
	}
}

func GetStringFromRole(r Role) string {
	switch r {
	case Trainee:
		return "trainee"
	case Admin:
		return "admin"
	default:
		return "trainee"
	}
}
