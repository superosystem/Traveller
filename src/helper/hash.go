package helper

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) string {
	defaultCost := bcrypt.DefaultCost

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), defaultCost)

	return string(hashedPassword)
}

func ComparePassword(hashedPassword string, password string) bool {
	ok := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return ok == nil
}
