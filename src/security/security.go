package security

import "golang.org/x/crypto/bcrypt"

// Hash receives a string(password) and generates a hash on it
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassowrd compares the hash with the password and checks if they match
func VerifyPassword(passwordHash, passwordString string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordString))
}
