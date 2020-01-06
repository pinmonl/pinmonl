package session

import (
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
)

// CookieStoreOpts defines the options of CookieStore initiation.
type CookieStoreOpts struct {
	Name     string
	MaxAge   time.Duration
	Secure   bool
	HashKey  []byte
	BlockKey []byte
}

// CookieStore uses cookie as session.
type CookieStore struct {
	sc     *securecookie.SecureCookie
	name   string
	maxAge int
	secure bool
}

var _ Store = &CookieStore{}

// NewCookieStore creates CookieStore with Opts.
func NewCookieStore(opts CookieStoreOpts) *CookieStore {
	sc := securecookie.New(opts.HashKey, opts.BlockKey)

	return &CookieStore{
		sc:     sc,
		name:   opts.Name,
		maxAge: int(opts.MaxAge),
		secure: opts.Secure,
	}
}

// Get retrieves Values from request.
func (s *CookieStore) Get(r *http.Request) (*Values, error) {
	cookie, err := r.Cookie(s.name)
	if err != nil {
		return nil, err
	}
	val := Values{}
	err = s.sc.Decode(s.name, cookie.Value, &val)
	if err != nil {
		return nil, err
	}
	return &val, nil
}

// Set returns cookie header to client.
func (s *CookieStore) Set(w http.ResponseWriter, val *Values) error {
	cookie, err := s.Cookie(val)
	if err != nil {
		return err
	}
	http.SetCookie(w, cookie)
	return nil
}

// Del removes the cookie from client.
func (s *CookieStore) Del(w http.ResponseWriter, r *http.Request) error {
	cookie, err := s.Cookie(nil)
	if err != nil {
		return err
	}

	cookie.MaxAge = 0
	http.SetCookie(w, cookie)
	return nil
}

// Cookie returns cookie with encoded value.
func (s *CookieStore) Cookie(val *Values) (*http.Cookie, error) {
	var enc string
	if val != nil {
		v, err := s.sc.Encode(s.name, val)
		if err != nil {
			return nil, err
		}
		enc = v
	}

	return &http.Cookie{
		Name:     s.name,
		Value:    enc,
		Path:     "/",
		MaxAge:   s.maxAge,
		Secure:   s.secure,
		HttpOnly: true,
	}, nil
}
