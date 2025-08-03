package security

import (
	"crypto/sha256"
	"encoding/hex"

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
    hash := hasher.Sum(nil)
    return hex.EncodeToString(hash), nil
}

func CompareToken(token, hash string) error {
    computedHash, err := HashToken(token)
    if err == nil && computedHash == hash {
        return nil
    }

    hasher := sha256.New()
    hasher.Write([]byte(token))
    preHashed := hasher.Sum(nil)
    return bcrypt.CompareHashAndPassword([]byte(hash), preHashed)
}