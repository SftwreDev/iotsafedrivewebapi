package utils

type Role string

const (
	super_admin Role = "super_admin"
	rescuer     Role = "rescuer"
	user        Role = "user"
)

func GetRole(r string) Role {
	return Role(r)
}
