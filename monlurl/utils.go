package monlurl

import (
	"golang.org/x/net/publicsuffix"
)

func TLDomain(host string) (string, error) {
	return publicsuffix.EffectiveTLDPlusOne(host)
}
