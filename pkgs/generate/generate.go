package generate

import (
	"math/rand"
	"time"
)

var (
	UserHashLength = 200

	CharsetAlphaNum = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789"
	CharsetAlphaNumSign = CharsetAlphaNum +
		"_-+$%()[]{}"
)

func StringWithCharset(n int, charset string) string {
	seed := rand.New(
		rand.NewSource(
			time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[seed.Intn(len(charset))]
	}
	return string(b)
}

func AlphaNum(n int) string {
	return StringWithCharset(n, CharsetAlphaNum)
}

func AlphaNumSign(n int) string {
	return StringWithCharset(n, CharsetAlphaNumSign)
}

func UserHash() string {
	return AlphaNum(UserHashLength)
}
