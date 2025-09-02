package tools

import "golang.org/x/crypto/bcrypt"

// Hash a string of password during sign-in or log-in.
func HashPassword(password string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(result), err
}

// Check and compare password while user is logging in.
func CompareHashPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
