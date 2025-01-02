package handler

import (
	"net/http"

	"github.com/axzilla/deeploy/internal/app/ui/pages"
)

func GetLanding(w http.ResponseWriter, r *http.Request) {
	pages.Dashboard().Render(r.Context(), w)
}
