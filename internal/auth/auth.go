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

// check admin token
// TODO: hash the superAdminToken for storage in memory and the requestToken
func ValidateSuperAdmin(superAdminToken string, requestToken string) bool {
	token1, err1 := base64.StdEncoding.DecodeString(superAdminToken)
	token2, err2 := base64.StdEncoding.DecodeString(requestToken)
	if err1 != nil || err2 != nil {
		return false
	}

	return hmac.Equal(token1, token2)
}

// === Token & Key Functions ===

// api key retrieval
// also used for super_admin_token
func GetAuthKeysValue(headers http.Header, prefix string) (string, error) {
	// value will look like:
	//   ApiKey <key string>
	if prefix == "" {
		prefix = "ApiKey"
	}

	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("header field 'Authorization' is absent")
	}

	keyString, ok := strings.CutPrefix(authHeader, prefix+" ")
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

// creates and returns a JWT
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
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

// returns a created refresh token
func MakeRefreshToken() (string, error) {
	data := make([]byte, 32)
	_, err := rand.Read(data)
	if err != nil {
		return "", err
	}

	secureString := hex.EncodeToString(data)
	return secureString, nil
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

// password is from a request, hash is from the db
func CheckPasswordHash(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Printf("Unable to compare hash and password: %s", err)
		return err
	}

	return nil
}
