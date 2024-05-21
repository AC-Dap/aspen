package auth

import "golang.org/x/crypto/bcrypt"

// hashPassword
// Hashes the given password using bcrypt's min cost
func hashPassword(password string) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Hash password with bcrypt's min cost
	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

// VerifyUser
// Verifies if we have a user with the given username and password
func VerifyUser(username string, password string, users []User) bool {
	passwordBytes := []byte(password)
	for _, user := range users {
		if user.Username != username {
			continue
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), passwordBytes); err == nil {
			return true
		}
	}
	return false
}

// VerifyUserToken
// Verifies if we have a user with the given username and token
func VerifyUserToken(username string, token string, users []User) bool {
	for _, user := range users {
		if user.Username == username && user.ActiveToken == token {
			return true
		}
	}
	return false
}
