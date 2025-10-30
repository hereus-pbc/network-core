package interfaces

type App interface {
	CheckPermission(permission string) bool
}
