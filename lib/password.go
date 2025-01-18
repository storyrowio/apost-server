package lib

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HashAndSalt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}

func ComparePassword(hashed string, current []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), current)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
