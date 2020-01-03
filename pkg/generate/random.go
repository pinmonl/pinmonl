package generate

import (
	"math/rand"
	"time"
)

// RandomString produces string with only alpha-numeric characters.
func RandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	return StringWithCharset(length, charset)
}

// StringWithCharset produces string from the given charset.
func StringWithCharset(length int, charset string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}
