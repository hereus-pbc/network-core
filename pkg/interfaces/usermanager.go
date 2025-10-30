package interfaces

type UserManager interface {
	Register(username string, firstName string, middleName string, lastName string, nameSuffix string, birthdate string, country string, isMale bool, timezone string, password string) error
	Get(username string) (User, error)
	Exists(username string) bool
}
