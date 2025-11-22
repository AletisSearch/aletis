package handlers

import (
	"net/http"

	"github.com/AletisSearch/aletis/web/templates"
	"github.com/AletisSearch/aletis/web/templates/home"
)

func Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "private")
		templates.Layout(home.Head(), home.Body()).Render(r.Context(), w)
	}
}
