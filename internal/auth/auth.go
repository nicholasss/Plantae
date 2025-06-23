package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// === Admin Token Functions ===

func ValidateSuperAdmin(superAdminToken string, requestToken string) bool {
// ValidateSuperAdmin is a cryptographicall secure function to check
// whether the token provided is the SuperAdminToken.
	token1, err1 := base64.StdEncoding.DecodeString(superAdminToken)
	token2, err2 := base64.StdEncoding.DecodeString(requestToken)
	if err1 != nil || err2 != nil {
		log.Print("Could not authenticate Super Admin.")
		return false
	}

	return hmac.Equal(token1, token2)
}

// === Token & Key Functions ===

func GetAuthKeysValue(headers http.Header, prefix string) (string, error) {
// GetAuthKeysValue returns the value in the 'Authorization' header of a request.
// Optionally provide a prefix to use before a token.
// i.e. "Bearer <token>", "ApiKey <token>", or "SuperAdminToken <token>"
// TODO: merge with function below
	// value will look like:
	//   ApiKey <key string>
	if prefix == "" {
		prefix = "ApiKey"
	}

	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("header field 'Authorization' is absent")
	}

	keyString, ok := strings.CutPrefix(authHeader, prefix+" ")
	if !ok {
		log.Printf("Unable to cut prefix off. Before: '%s' After: '%s'", authHeader, keyString)
		return "", errors.New("unable to find key in headers")
	}

	return keyString, nil
}

func GetBearerToken(headers http.Header) (string, error) {
// GetBearerToken returns the access token from a requests headers.
// TODO: merge with function above
	// value will look like
	//   Bearer <token_string>

	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("header field 'Authorization' is absent")
	}

	tokenString, ok := strings.CutPrefix(authHeader, "Bearer ")
	if !ok {
		if strings.Contains(authHeader, "SuperAdminToken") {
			log.Println("Super-admin token supplied where admin's access token is required.")
			return "", errors.New("super-admin token provided, please provide admin's access token instead")
		}

		log.Printf("Unable to cut prefix off. Before: '%s' After: '%s'\n", authHeader, tokenString)
		return "", errors.New("unable to find token in headers")
	}

	// log.Printf("Returned the JWT successfuly from headers.\n")
	return tokenString, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
// MakeJWT provides a fresh access token to a particular user for a given duration.
	currentTime := time.Now().UTC()
	expirationTime := currentTime.UTC().Add(expiresIn)

	signingMethod := jwt.SigningMethodHS256
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(currentTime),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(signingMethod, claims)

	// HMAC signing method requires the type []byte
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		log.Printf("Error signing JWT: %s", err)
		return "", err
	}

	// log.Printf("Generated new token: %s", signedToken)
	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
// ValidateJWT checks a users access token and ensures that it is valid.
// It will return a user id (uuid) when successful.
	claims := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return uuid.Nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(tokenSecret), nil
		})
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, err
	}

	return userUUID, nil
}

func MakeRefreshToken() (string, error) {
// MakeRefreshToken provides a fresh refresh token.
	data := make([]byte, 32)
	_, err := rand.Read(data)
	if err != nil {
		return "", err
	}

	secureString := hex.EncodeToString(data)
	return secureString, nil
}

// === User Password Functions ===

func HashPassword(rawPassword string) (string, error) {
// HashPassword takes a raw password and returns a hashed version, utilizing bcrypt.
	if rawPassword == "" {
		log.Print("Empty password provided.")
		return "", errors.New("unable to hash empty password")
	}

	rawPasswordData := []byte(rawPassword)
	rawPassword = "" // GC collection
	if len(rawPasswordData) > 72 {
		return "", errors.New("unable to hash password longer than 72 bytes")
	}

	hashedPasswordData, err := bcrypt.GenerateFromPassword(rawPasswordData, bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Unable to hash password: %s", err)
		return "", err
	}
	rawPasswordData = nil // GC collection

	return string(hashedPasswordData), nil
}

func CheckPasswordHash(password, hashedPassword string) error {
// CheckPasswordHash takes a raw password from a client, and a hashed password from the server.
// If the hased raw password (from client) does not match the stored hashed password (from database),
// it will return an error.
// If they match, then it will return nil.
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Printf("Unable to compare hash and password: %s", err)
		return err
	}

	return nil
}
