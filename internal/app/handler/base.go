package handler

import (
	"net/http"

	"github.com/axzilla/deeploy/internal/app/ui/pages"
)

type BaseHandler struct{}

func NewBaseHandler() BaseHandler {
	return BaseHandler{}
}

func (*BaseHandler) GetLanding(w http.ResponseWriter, r *http.Request) {
	pages.Dashboard().Render(r.Context(), w)
}
