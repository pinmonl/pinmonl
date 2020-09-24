package web

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/markbates/pkger"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/exchange"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/pubsub"
	"github.com/pinmonl/pinmonl/queue"
	"github.com/pinmonl/pinmonl/store"
)

type Server struct {
	Txer        database.Txer
	TokenSecret []byte
	TokenExpire time.Duration
	TokenIssuer string
	Queue       *queue.Manager
	Exchange    *exchange.Manager
	Pubsub      pubsub.Pubsuber

	ExchangeEnabled bool
	DefaultUserID   string
	DevServer       string

	Images    *store.Images
	Monls     *store.Monls
	Monpkgs   *store.Monpkgs
	Pinls     *store.Pinls
	Pinpkgs   *store.Pinpkgs
	Pkgs      *store.Pkgs
	Sharepins *store.Sharepins
	Shares    *store.Shares
	Sharetags *store.Sharetags
	Stats     *store.Stats
	Taggables *store.Taggables
	Tags      *store.Tags
	Users     *store.Users
}

func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(s.authenticate())
	r.Mount("/api", s.APIRouter())
	r.Mount("/", s.WebRouter())
	return r
}

func (s *Server) APIRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	if !s.hasDefaultUser() {
		r.Post("/signup", s.signupHandler)
		r.Post("/login", s.loginHandler)
	}
	r.With(s.authorize()).
		Post("/refresh", s.refreshHandler)

	r.Route("/pin", func(r chi.Router) {
		r.Use(s.authorize())
		r.With(s.pagination()).
			Get("/", s.pinlListHandler)
		r.Post("/", s.pinlCreateHandler)
		r.Route("/{pinl}", func(r chi.Router) {
			r.Use(s.bindPinl())
			r.Get("/", s.pinlHandler)
			r.Put("/", s.pinlUpdateHandler)
			r.Delete("/", s.pinlDeleteHandler)
			r.Post("/image", s.pinlUploadImageHandler)
		})
	})

	r.Get("/card", s.fetchCardHandler)

	r.Route("/tag", func(r chi.Router) {
		r.Use(s.authorize())
		r.With(s.pagination()).
			Get("/", s.tagListHandler)
		r.Post("/", s.tagCreateHandler)
		r.Route("/{tag}", func(r chi.Router) {
			r.Use(s.bindTag())
			r.Get("/", s.tagHandler)
			r.Put("/", s.tagUpdateHandler)
			r.Delete("/", s.tagDeleteHandler)
		})
	})

	r.Route("/pkg", func(r chi.Router) {
		r.With(s.pagination()).
			Get("/", s.pkgListHandler)
		r.Post("/", s.pkgCreateHandler)
	})

	r.Route("/stat", func(r chi.Router) {
		r.With(s.pagination()).
			Get("/", s.statListHandler)
	})

	return r
}

func (s *Server) WebRouter() chi.Router {
	r := chi.NewRouter()
	r.With(s.bindImage()).Get("/image/{image}", s.imageHandler)
	r.Handle("/*", s.webHandler())
	return r
}

func (s *Server) pagination() func(http.Handler) http.Handler {
	return request.Pagination("page", "pageSize", 10)
}

func (s *Server) hasDefaultUser() bool {
	return s.DefaultUserID != ""
}

func (s *Server) defaultUser() *model.User {
	user, _ := s.Users.Find(context.TODO(), s.DefaultUserID)
	return user
}

func (s *Server) webHandler() http.Handler {
	var handler http.Handler
	if s.DevServer != "" {
		u, err := url.Parse(s.DevServer)
		if err != nil {
			return nil
		}
		handler = httputil.NewSingleHostReverseProxy(u)
	} else {
		pkgdir, _ := pkger.Open("/webui/dist")
		handler = http.FileServer(pkgdir)
	}

	indexHandler := func(w http.ResponseWriter, r *http.Request) {
		indexTmpl, err := getIndexTemplate(handler)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		baseURL := &url.URL{
			Scheme: "http",
			Host:   r.Host,
		}

		data := map[string]interface{}{
			"BaseURL":         baseURL.String(),
			"BasePrefix":      "",
			"HasDefaultUser":  s.hasDefaultUser(),
			"ExchangeEnabled": s.ExchangeEnabled,
		}

		indexTmpl.Execute(w, data)
		return
	}

	r := chi.NewRouter()
	r.Get("/", indexHandler)
	r.Handle("/*", handler)

	return r
}

func getIndexTemplate(h http.Handler) (*template.Template, error) {
	w := newLazyResponseWriter()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		return nil, err
	}

	h.ServeHTTP(w, r)

	tmpl, err := template.New("").Parse(w.content.String())
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

type lazyResponseWriter struct {
	code    int
	content *bytes.Buffer
	header  http.Header
}

func newLazyResponseWriter() *lazyResponseWriter {
	return &lazyResponseWriter{
		content: &bytes.Buffer{},
		header:  http.Header{},
	}
}

func (lw *lazyResponseWriter) Write(bs []byte) (int, error) {
	return lw.content.Write(bs)
}

func (lw *lazyResponseWriter) WriteHeader(code int) {
	lw.code = code
}

func (lw *lazyResponseWriter) Header() http.Header {
	return lw.header
}
