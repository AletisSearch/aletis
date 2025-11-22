package aletis

import (
	"context"
	"log/slog"
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

	aiClient := aiclient.NewClient(conf.OpenAIURL, conf.OpenAIKey, q)

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
	r.Get("/", handlers.Home())
	r.Route("/search", func(r chi.Router) {
		if conf.Public {
			r.Use(httprate.LimitByRealIP(10, time.Minute))
		}
		r.Get("/", handlers.Search(aiClient, searchClient))
	})
	r.Get("/icons/{domain}", handlers.Icons(q))
	r.Handle("/assets/*", handlers.Assets(conf.Dev))
	return r, nil
}
