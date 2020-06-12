package github

import (
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	globalTokens = NewTokenStore(nil)
)

var (
	ErrNoToken         = errors.New("github: no token available")
	ErrTokenReachLimit = errors.New("github: token reaches the limit")
	ErrTokenUsing      = errors.New("github: token is using")
)

type TokenStore struct {
	infos []*TokenInfo
}

func NewTokenStore(tokens []string) *TokenStore {
	t := &TokenStore{}
	t.Add(tokens)
	return t
}

func (t *TokenStore) Add(tokens []string) {
	for i := range tokens {
		t.infos = append(t.infos, NewTokenInfo(tokens[i]))
	}
}

func (t *TokenStore) Get() (*TokenInfo, error) {
	for _, info := range t.infos {
		if info.using {
			continue
		}
		if err := info.Valid(); err == nil {
			return info, nil
		}
	}
	return nil, ErrNoToken
}

func (t *TokenStore) Set(tokens []string) {
	t.infos = make([]*TokenInfo, 0)
	t.Add(tokens)
}

func AddToken(tokens []string) {
	globalTokens.Add(tokens)
}

func GetToken() (*TokenInfo, error) {
	return globalTokens.Get()
}

func SetToken(tokens []string) {
	globalTokens.Set(tokens)
}

type TokenInfo struct {
	*sync.Mutex
	token     string
	remaining int
	reset     time.Time
	using     bool
}

func NewTokenInfo(token string) *TokenInfo {
	return &TokenInfo{
		Mutex: &sync.Mutex{},
		token: token,
	}
}

func (t *TokenInfo) Lock() {
	t.using = true
	t.Mutex.Lock()
}

func (t *TokenInfo) Unlock() {
	t.Mutex.Unlock()
	t.using = false
}

func (t *TokenInfo) Valid() error {
	if t.remaining == 0 && time.Now().Before(t.reset) {
		return ErrTokenReachLimit
	}
	return nil
}

func (t *TokenInfo) UpdateFromHeader(header http.Header) error {
	var (
		remaining int
		reset     time.Time
	)
	if r, err := strconv.Atoi(header.Get("X-RateLimit-Remaining")); err == nil {
		remaining = r
	} else {
		return err
	}
	if r, err := strconv.ParseInt(header.Get("X-RateLimit-Reset"), 10, 64); err == nil {
		reset = time.Unix(r, 0)
	} else {
		return err
	}

	t.remaining = remaining
	t.reset = reset
	return nil
}
