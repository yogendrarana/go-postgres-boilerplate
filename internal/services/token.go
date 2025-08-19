package services

import (
	"fmt"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 10
)

func GenerateAccessToken(userID uuid.UUID) (string, error) {
	accessExpiry := time.Now().Add(time.Minute * 15).Unix()
	claims := jwt.MapClaims{"exp": accessExpiry, "sub": userID}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tkn.SignedString([]byte(os.Getenv("ACCESS_JWT_SECRET")))
}

func ValidateAccessToken(accessToken string) (bool, *uuid.UUID) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return false, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, nil
	}

	// sub is stored as string if uuid, or could parse
	var id uuid.UUID
	switch v := claims["sub"].(type) {
	case string:
		parsed, parseErr := uuid.Parse(v)
		if parseErr != nil {
			return false, nil
		}
		id = parsed
	case fmt.Stringer:
		parsed, parseErr := uuid.Parse(v.String())
		if parseErr != nil {
			return false, nil
		}
		id = parsed
	default:
		return false, nil
	}
	return true, &id
}

func GenerateRefreshTokenAndHash() (string, string, error) {
	refreshToken := uuid.New().String()
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcryptCost)
	if err != nil {
		return "", "", err
	}
	return refreshToken, string(hashedToken), nil
}
