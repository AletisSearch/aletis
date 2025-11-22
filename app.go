package aletis

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	aiclient "github.com/AletisSearch/aletis/internal/aiClient"
	"github.com/AletisSearch/aletis/internal/db"
	"github.com/AletisSearch/aletis/internal/handlers"
	"github.com/AletisSearch/aletis/internal/searxng"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func NewApp(ctx context.Context, wg *sync.WaitGroup, conf *Config, q *db.Queries) (*chi.Mux, error) {
	var aiClient *aiclient.Client
	if conf.AIEnabled {
		aiClient = aiclient.NewClient(conf.OpenAIURL, conf.OpenAIKey, q)
	}
	searchClient := searxng.NewClient(conf.SearxngHost, q)

	wg.Go(func() {
		<-ctx.Done()
		slog.Info("Closing Search Client")
		searchClient.Close()
	})

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)
	r.Use(middleware.SetHeader("X-Frame-Options", "DENY"))
	r.Use(middleware.SetHeader("X-Content-Type-Options", "nosniff"))
	r.Use(middleware.SetHeader("Referrer-Policy", "strict-origin-when-cross-origin"))
	//r.Use(middleware.SetHeader("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload"))
	r.Use(middleware.SetHeader("Cross-Origin-Opener-Policy", "same-origin"))
	r.Use(middleware.SetHeader("Cross-Origin-Embedder-Policy", "require-corp"))
	r.Use(middleware.SetHeader("Cross-Origin-Resource-Policy", "same-site"))
	r.Use(middleware.SetHeader("Permissions-Policy", "geolocation=(), camera=(), microphone=(), interest-cohort=()"))

	r.Get("/", handlers.Home())
	r.Route("/search", func(r chi.Router) {
		if conf.Public {
			r.Use(httprate.LimitByRealIP(10, time.Minute))
		}
		r.Get("/", handlers.Search(aiClient, searchClient))
	})
	r.Get("/icons/{domain}", handlers.Icons(q))
	r.Handle("/assets/*", handlers.Assets(conf.Dev))

	r.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=3600")
		w.Write([]byte(`User-agent: *
Disallow: /search
Disallow: /icons
Disallow: /assets`))
	})
	return r, nil
}
