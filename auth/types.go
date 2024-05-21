package auth

type User struct {
	Username     string
	PasswordHash string
	ActiveToken  string
}

type AuthTOML struct {
	Users []User
}
