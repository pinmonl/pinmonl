package middleware

import (
	"net/http"

	"github.com/pinmonl/pinmonl/logx"
	"github.com/pinmonl/pinmonl/store"
)

// EnableTransaction begins transaction per HTTP request and commits when response is made.
func EnableTransaction(s store.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := NewTxWriter(w)

			ctx, err := s.BeginTx(r.Context())
			if err != nil {
				logx.Fatalf("api: fails to start transaction, err: %s", err)
			}
			next.ServeHTTP(ww, r.WithContext(ctx))

			if ww.Fails() {
				err = s.Rollback(ctx)
				if err != nil {
					logx.Fatalf("api: fails to rollback transaction, err: %s", err)
				}
			} else {
				err = s.Commit(ctx)
				if err != nil {
					logx.Fatalf("api: fails to commit transaction, err: %s", err)
				}
			}
		}
		return http.HandlerFunc(fn)
	}
}

// TxWriter wraps http.ResponseWriter for status code handler.
type TxWriter interface {
	http.ResponseWriter
	Fails() bool
}

// NewTxWriter creates tx writer.
func NewTxWriter(w http.ResponseWriter) TxWriter {
	bw := &basicWriter{w: w}
	return bw
}

type basicWriter struct {
	w    http.ResponseWriter
	code int
}

// Header returns response header.
func (b *basicWriter) Header() http.Header {
	return b.w.Header()
}

// Write writes to response body
// and sets status code to 200 if has not set yet.
func (b *basicWriter) Write(body []byte) (int, error) {
	l, err := b.w.Write(body)
	if err == nil && b.code == 0 {
		b.code = 200
	}
	return l, err
}

// WriteHeader writes status code.
func (b *basicWriter) WriteHeader(code int) {
	if b.code == 0 {
		b.code = code
	}
	b.w.WriteHeader(code)
}

// Fails reports failed status code.
func (b *basicWriter) Fails() bool {
	switch c := b.code; {
	case c < 200:
		return true
	case c >= 400:
		return true
	default:
		return false
	}
}
