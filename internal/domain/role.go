package domain

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleAuthor   Role = "author"
	RoleReviewer Role = "reviewer"
	RoleObserver Role = "observer"
)

func ParseRole(s string) Role {
	switch s {
	case string(RoleAdmin):
		return RoleAdmin
	case string(RoleAuthor):
		return RoleAuthor
	case string(RoleReviewer):
		return RoleReviewer
	case string(RoleObserver):
		return RoleObserver
	default:
		return RoleObserver
	}
}
