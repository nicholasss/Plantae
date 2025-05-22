package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// === Admin Token Functions ===

// check admin token
func authorizeSuperAdmin(superAdminToken string, requestToken string) bool {
	result := strings.Compare(superAdminToken, requestToken)
	return result == 0 // true if same, false if not
}

// === Token & Key Functions ===

// api key retrieval
// also used for super_admin_token
func GetAPIKey(headers http.Header) (string, error) {
	// value will look like:
	//   ApiKey <key string>

	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("header field 'authorization' is absent")
	}

	keyString, ok := strings.CutPrefix(authHeader, "ApiKey ")
	if !ok {
		log.Printf("Unable to cut prefix off. Before: '%s' After: '%s'", authHeader, keyString)
		return "", errors.New("unable to find key in headers")
	}

	return keyString, nil
}

// user access token or refresh token retrieval
func GetBearerToken(headers http.Header) (string, error) {
	// value will look like
	//   Bearer <token_string>

	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("header field 'authorization' is absent")
	}

	tokenString, ok := strings.CutPrefix(authHeader, "Bearer ")
	if !ok {
		log.Printf("Unable to cut prefix off. Before: '%s' After: '%s'\n", authHeader, tokenString)
		return "", errors.New("unable to find token in headers")
	}

	log.Printf("Returned the JWT successfuly from headers.\n")
	return tokenString, nil
}

// === User Password Functions ===

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
