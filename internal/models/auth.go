package models

type Role string

func (r Role) String() string {
	return string(r)
}

const (
	ADMIN = Role("admin")
	USER = Role("user")
	VIEWER = Role("viewer")
)

type Permission string

func (p Permission) String() string {
	return string(p)
}

const (
	CREATE = Permission("create")
	UPDATE = Permission("update")
	DELETE = Permission("delete")
	VIEW   = Permission("view")
)

var roles = map[Role][]Permission {
	ADMIN: {CREATE, UPDATE, DELETE, VIEW},
	USER:   {CREATE, VIEW},
	VIEWER: {VIEW},
}

func CheckPermission(role string) ([]Permission, bool) {
	p, ok := roles[Role(role)]
	return p, ok
}

func GetRolesList() []Role {
	keys := make([]Role, len(roles))
	for k := range roles {
		keys = append(keys, k)
	}
	return keys
}