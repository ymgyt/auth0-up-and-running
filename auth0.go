package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ymgyt/appkit/handlers"
	"go.uber.org/zap"
)

// Auth0 -
type Auth0 struct {
	ts     *handlers.TemplateSet
	logger *zap.Logger
	cfg    *config
}

// RenderLogin -
func (a *Auth0) RenderLogin(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	view := struct {
		Domain   string
		ClientID string
	}{
		Domain:   a.cfg.Auth0Domain,
		ClientID: a.cfg.Auth0ClientID,
	}
	if err := a.ts.ExecuteTemplate(w, "login", &view); err != nil {
		a.logger.Error("render_login", zap.Error(err))
	}
}
