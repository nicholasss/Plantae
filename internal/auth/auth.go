/*
Package auth provides cryptographically secure functions for passwords or javascript web tokens.
*/
package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// === Admin Token Functions ===

// ValidateSuperAdmin is a cryptographicall secure function to check
// whether the token provided is the SuperAdminToken.
func ValidateSuperAdmin(superAdminToken string, requestToken string, sl *slog.Logger) bool {
	token1, err1 := base64.StdEncoding.DecodeString(superAdminToken)
	token2, err2 := base64.StdEncoding.DecodeString(requestToken)
	if err1 != nil || err2 != nil {
		sl.Debug("One of the provided tokens is empty")
		return false
	}

	return hmac.Equal(token1, token2)
}

// === Token & Key Functions ===

// fetches the authorization header
// is intended to be used by other functions
func getAuthHeader(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("header field 'Authorization' is absent")
	}

	return authHeader, nil
}

// cuts the prefix from the token
// returning just the token itsefl
func removeTokenPrefix(prefixedToken, prefix string) (string, error) {
	token, ok := strings.CutPrefix(prefixedToken, prefix+" ")
	if !ok {
		return "", errors.New("unable to find key in headers")
	}

	return token, nil
}

// GetSuperAdminToken returns the SuperAdminToken from headers
// -- SuperAdminToken <token_string>
func GetSuperAdminToken(headers http.Header, sl *slog.Logger) (string, error) {
	authValue, err := getAuthHeader(headers)
	if err != nil {
		return "", err
	}

	token, err := removeTokenPrefix(authValue, "SuperAdminToken")
	if err != nil {
		sl.Debug("Unable to cut prefix from Authorization header")
		return "", err
	}

	return token, nil
}

// GetBearerToken returns the Bearer token from headers
// -- Bearer <token_string>
func GetBearerToken(headers http.Header, sl *slog.Logger) (string, error) {
	authValue, err := getAuthHeader(headers)
	if err != nil {
		return "", err
	}

	token, err := removeTokenPrefix(authValue, "Bearer")
	if err != nil {
		sl.Debug("Unable to cut prefix from Authorization header")
		return "", err
	}

	return token, nil
}

// MakeJWT provides a fresh access token to a particular user for a given duration.
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration, sl *slog.Logger) (string, error) {
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
		sl.Debug("Unable to sign JWT", "error", err)
		return "", err
	}

	// sl.Debug("Generated new token", "raw access token", signedToken)
	return signedToken, nil
}

// ValidateJWT checks a users access token and ensures that it is valid.
// It will return a user id (uuid) when successful.
func ValidateJWT(tokenString, tokenSecret string, sl *slog.Logger) (uuid.UUID, error) {
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

// MakeRefreshToken provides a fresh refresh token.
func MakeRefreshToken(sl *slog.Logger) (string, error) {
	data := make([]byte, 32)
	_, err := rand.Read(data)
	if err != nil {
		sl.Debug("Unable to read random data", "error", err)
		return "", err
	}

	secureString := hex.EncodeToString(data)
	return secureString, nil
}

// === User Password Functions ===

// HashPassword takes a raw password and returns a hashed version, utilizing bcrypt.
func HashPassword(rawPassword string, sl *slog.Logger) (string, error) {
	if rawPassword == "" {
		sl.Debug("Empty password was passed in")
		return "", errors.New("unable to hash empty password")
	}

	rawPasswordData := []byte(rawPassword)
	rawPassword = "" // GC collection
	if len(rawPasswordData) > 72 {
		return "", errors.New("unable to hash password longer than 72 bytes")
	}

	hashedPasswordData, err := bcrypt.GenerateFromPassword(rawPasswordData, bcrypt.DefaultCost)
	if err != nil {
		sl.Debug("Unable to hash password", "error", err)
		return "", err
	}
	rawPasswordData = nil // GC collection

	return string(hashedPasswordData), nil
}

// CheckPasswordHash takes a raw password from a client, and a hashed password from the server.
// If the hased raw password (from client) does not match the stored hashed password (from database),
// it will return an error.
// If they match, then it will return nil.
func CheckPasswordHash(password, hashedPassword string, sl *slog.Logger) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		sl.Debug("Unable to compare provided raw password and hashed password", "error", err)
		return err
	}

	return nil
}
