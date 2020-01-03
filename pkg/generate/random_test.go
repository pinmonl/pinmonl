package generate

import (
	"regexp"
	"testing"
)

func TestRandomString(t *testing.T) {
	lens := []int{10, 20, 30, 40, 50, 60}
	pattern := regexp.MustCompile("^[a-zA-Z0-9]+$")

	for _, l := range lens {
		got := RandomString(l)
		match := pattern.MatchString(got)
		if len(got) != l || !match {
			t.Errorf("random string generation fails, got: %s", got)
		}
	}
}
