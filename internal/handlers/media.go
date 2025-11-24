package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/AletisSearch/aletis/internal/db"
	"github.com/AletisSearch/aletis/internal/icons"
	"github.com/AletisSearch/aletis/web"
	"github.com/go-playground/validator/v10"
)

func Assets(dev bool) http.Handler {
	return http.StripPrefix("/assets/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if dev {
			p, err := web.GetFilepath(r.URL.Path)
			if err != nil {
				if !errors.Is(err, web.ErrFileNotFound) {
					slog.Error("Err getting filepath", "ERR", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
			rp, err := web.GetFilepath(r.URL.RawPath)
			if err != nil {
				if !errors.Is(err, web.ErrFileNotFound) {
					slog.Error("Err getting filepath", "ERR", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
			if p != "" || rp != "" {
				r2 := new(http.Request)
				*r2 = *r
				r2.URL = new(url.URL)
				*r2.URL = *r.URL
				r2.URL.Path = p
				r2.URL.RawPath = rp
				w.Header().Set("Cache-Control", "no-cache")
				http.StripPrefix("/assets/", http.FileServer(http.FS(web.AssetsFs))).ServeHTTP(w, r2)
				return
			}
		}
		w.Header().Add("Cache-Control", "public, max-age=2629800, immutable")
		http.FileServer(http.FS(web.AssetsFs)).ServeHTTP(w, r)
	}))
}

func Icons(db *db.Queries) http.HandlerFunc {
	c := icons.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		domainRaw := r.PathValue("domain")
		i, err := c.Get(r.Context(), domainRaw)
		if err != nil {
			if i == nil {
				var validateErrs validator.ValidationErrors
				if errors.As(err, &validateErrs) {
					w.WriteHeader(http.StatusBadRequest)
					slog.Error("bad icon request", "ERROR", err)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				slog.Error("something went wrong with icon request", "ERROR", err)
				return
			}
			slog.Error("able to get icon but errored", "ERROR", err)
		}
		w.Header().Add("Content-Type", i.ContentType)
		w.Header().Add("Cache-Control", "public, max-age="+strconv.Itoa(int(i.Expiration.Seconds())))
		w.Write(i.IconBytes)
	}
}
