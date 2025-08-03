package security

import (
	"crypto/sha256"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func ComparePasswords(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func HashToken(token string) (string, error) {
    hasher := sha256.New()
    hasher.Write([]byte(token))
    preHashed := hasher.Sum(nil)

    hash, err := bcrypt.GenerateFromPassword(preHashed, bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

func CompareToken(token, hash string) error {
    hasher := sha256.New()
    hasher.Write([]byte(token))
    preHashed := hasher.Sum(nil)

    return bcrypt.CompareHashAndPassword([]byte(hash), preHashed)
}