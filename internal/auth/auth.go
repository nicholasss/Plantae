package auth

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// hashes password using the bcrypt golang library
func HashPassword(rawPassword string) (string, error) {
	if rawPassword == "" {
		log.Print("Empty password provided.")
		return "", fmt.Errorf("unable to hash empty password")
	}

	rawPasswordData := []byte(rawPassword)
	rawPassword = "" // GC collection
	if len(rawPasswordData) > 72 {
		return "", fmt.Errorf("unable to hash password longer than 72 bytes")
	}

	hashedPasswordData, err := bcrypt.GenerateFromPassword(rawPasswordData, bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Unable to hash password: %s", err)
		return "", err
	}
	rawPasswordData = nil // GC collection

	return string(hashedPasswordData), nil
}
