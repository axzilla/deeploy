package handler

import (
	"net/http"

	"github.com/axzilla/deeploy/internal/app/ui/pages"
)

type BaseHandler struct{}

func NewBaseHandler() BaseHandler {
	return BaseHandler{}
}

func (*BaseHandler) DashboardView(w http.ResponseWriter, r *http.Request) {
	pages.Dashboard().Render(r.Context(), w)
}

func (*BaseHandler) LandingView(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
