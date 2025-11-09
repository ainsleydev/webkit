package random

import (
	"fmt"
	"math/rand"
	"time"
)

// String generates a random hexadecimal string of the specified length.
func String(length int64) string {
	rand.Seed(time.Now().UnixNano()) //nolint
	b := make([]byte, length)
	rand.Read(b) //nolint
	return fmt.Sprintf("%x", b)[:length]
}

// Alpha generates a random alphabetic string of the specified length.
func Alpha(length int64) string {
	return generate(length, false)
}

// AlphaNum generates a random alphanumeric string of the specified length.
func AlphaNum(length int64) string {
	return generate(length, false)
}

// Seq generates a random alphabetic sequence of the specified length.
func Seq(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano()) //nolint
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// generate is a helper function for generating random character runs
// with no special characters
func generate(n int64, numeric bool) string {
	characterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	if !numeric {
		characterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	}
	b := make([]rune, n)
	for i := range b {
		b[i] = characterRunes[rand.Intn(len(characterRunes))]
	}
	return string(b)
}
