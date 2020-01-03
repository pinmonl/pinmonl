package response

import "net/http"

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
