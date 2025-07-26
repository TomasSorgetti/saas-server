package security

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

type VerificationClaims struct {
	UserID                int64     `json:"user_id"`
	Email                 string    `json:"email"`
	VerificationExpiresAt time.Time `json:"verification_expires_at"`
	jwt.RegisteredClaims
}

func CreateAccessToken(userID int64) (string, error) {
	secret := os.Getenv("JWT_ACCESS_SECRET")
	if secret == "" {
		return "", errors.New("JWT_ACCESS_SECRET not set")
	}

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func CreateRefreshToken(userID int64) (string, error) {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		return "", errors.New("JWT_REFRESH_SECRET not set")
	}

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func CreateVerificationToken(userID int64, email string, verificationExpiresAt time.Time) (string, error) {
	secret := os.Getenv("JWT_VERIFICATION_SECRET") 
	if secret == "" {
		return "", errors.New("JWT_VERIFICATION_SECRET not set")
	}

	claims := &VerificationClaims{
		UserID:                userID,
		Email:                 email, 
		VerificationExpiresAt: verificationExpiresAt,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateAccessToken(tokenStr string) (int64, error) {
	return validateToken(tokenStr, os.Getenv("JWT_ACCESS_SECRET"))
}

func ValidateRefreshToken(tokenStr string) (int64, error) {
	return validateToken(tokenStr, os.Getenv("JWT_REFRESH_SECRET"))
}

func ValidateVerificationToken(tokenStr string) (int64, string, time.Time, error) {
	secret := os.Getenv("JWT_VERIFICATION_SECRET")
	if secret == "" {
		return 0, "", time.Time{}, errors.New("JWT_VERIFICATION_SECRET not set")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &VerificationClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return 0, "", time.Time{}, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*VerificationClaims)
	if !ok {
		return 0, "", time.Time{}, errors.New("invalid claims")
	}

	return claims.UserID, claims.Email, claims.VerificationExpiresAt, nil
}

func validateToken(tokenStr, secret string) (int64, error) {
	if secret == "" {
		return 0, errors.New("JWT secret not set")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}

	return 0, ErrInvalidToken
}
