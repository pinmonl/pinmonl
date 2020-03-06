package web

import (
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/webui"
)

// Handler returns the router.
func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()

	if s.devSvr != "" {
		r.HandleFunc("/*", handleWebuiDevServer(s.devSvr))
	} else {
		r.HandleFunc("/*", handleWebuiStatic())
	}

	return r
}

func handleWebuiStatic() http.HandlerFunc {
	box := webui.PackrBox()
	index, _ := box.Find("index.html")
	return func(w http.ResponseWriter, r *http.Request) {
		fp := r.URL.String()
		c, err := box.Find(fp)
		if err != nil {
			w.Write(index)
			return
		}

		ct := mime.TypeByExtension(filepath.Ext(fp))
		if ct == "" {
			ct = http.DetectContentType(c[:512])
		}
		w.Header().Add("Content-Type", ct)
		w.Write(c)
	}
}

func handleWebuiDevServer(devSvr string) http.HandlerFunc {
	url, _ := url.Parse(devSvr)
	devRP := httputil.NewSingleHostReverseProxy(url)
	return func(w http.ResponseWriter, r *http.Request) {
		devRP.ServeHTTP(w, r)
	}
}
