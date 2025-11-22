package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	aiclient "github.com/AletisSearch/aletis/internal/aiClient"
	"github.com/AletisSearch/aletis/internal/searxng"
	"github.com/AletisSearch/aletis/web/templates"
	"github.com/AletisSearch/aletis/web/templates/search"
	"github.com/a-h/templ"
)

type sd struct {
	Title string
	URL   string
	MD    string
}

func Search(aiClient *aiclient.Client, searchClient *searxng.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := strings.TrimSpace(r.URL.Query().Get("q"))
		if query == "" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		queryWSpaces := strings.ReplaceAll(query, "+", " ")

		dataChan := make(chan templ.Component)
		var wg sync.WaitGroup

		wg.Go(func() {
			sr, err := searchClient.Search(r.Context(), query)
			if err != nil {
				if sr == nil {
					slog.Error("unable to get searxng response", "ERROR", err)
					dataChan <- search.R("result", "Something went wrong")
					return
				}
				slog.Error("able to get searxng response but errored", "ERROR", err)
			}

			if len(sr.Results) == 0 {
				dataChan <- search.R("result", "No results")
				return
			}

			dataChan <- search.Results(sr)
		})
		if aiClient != nil {
			wg.Go(func() {
				data, err := aiClient.RunQueryExpand(r.Context(), fmt.Sprintf("[%s]", queryWSpaces))
				if err != nil {
					slog.Error("unable to get ai recommendations", "ERROR", err)
					dataChan <- search.R("recommendations", "Something went wrong")
					return
				}

				dataChan <- search.Recommendations(strings.Split(data.Content, "\n"))
			})
		}
		go func() {
			wg.Wait()
			close(dataChan)
		}()

		aiEnabled := aiClient != nil
		c := templates.Layout(search.Head(), search.Body(queryWSpaces, aiEnabled, dataChan))

		w.Header().Add("Cache-Control", "private")
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		templ.Handler(c, templ.WithStreaming()).ServeHTTP(w, r)
	}
}
