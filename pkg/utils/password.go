package utils

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
	"wtm-backend/pkg/logger"
)

func ComparePassword(ctx context.Context, passCrypt, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passCrypt), []byte(password))
	if err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		logger.Error(ctx, "Error comparing password:", err.Error())
		return false
	}
	return err == nil
}

func GeneratePassword(ctx context.Context, password string) (string, error) {
	passCrypt, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(ctx, "Error generating password", err.Error())
		return "", err
	}

	return string(passCrypt), nil
}

func GenerateSafeRandomString(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// ASCII printable range: 33 ('!') to 126 ('~')
	// Tapi kita exclude: ' " ( ) { } [ ]
	exclude := map[byte]bool{
		'\'': true, '"': true,
		'(': true, ')': true,
		'{': true, '}': true,
		'[': true, ']': true,
	}

	result := make([]byte, 0, n)
	for len(result) < n {
		c := byte(r.Intn(126-33+1) + 33)
		if !exclude[c] {
			result = append(result, c)
		}
	}
	return string(result)
}
