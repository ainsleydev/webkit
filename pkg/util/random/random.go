package random

import "math/rand"

// String generates a random of characters only (upper or lowercase)
// with no numbers based on the input length provided.
func String(length int64) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

// Alpha generates a random string (with no special characters)
// of n length with no numeric characters.
func Alpha(length int64) string {
	return generate(length, false)
}

// AlphaNum generates a random string (with no special characters)
// of n length with numeric characters.
func AlphaNum(length int64) string {
	return generate(length, false)
}

// Seq generates a random seeded sequence by the number of n
// input.
func Seq(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// generate is a helper function for generating random character runs
// with no special characters
func generate(n int64, numeric bool) string {
	var characterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	if !numeric {
		characterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	}
	b := make([]rune, n)
	for i := range b {
		b[i] = characterRunes[rand.Intn(len(characterRunes))]
	}
	return string(b)
}
