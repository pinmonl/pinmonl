package exchange

import (
	"net/http"

	"github.com/pinmonl/pinmonl/pinmonl-go"
)

type Transport struct {
	token string
	base  http.RoundTripper
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}

	if t.token != "" {
		req.Header.Add("Authorization", "Bearer "+t.token)
	}
	return base.RoundTrip(req)
}

func newPMClient(addr, token string) *pinmonl.Client {
	tp := &Transport{token: token}
	client := &http.Client{Transport: tp}
	return pinmonl.NewClient(addr, client)
}
