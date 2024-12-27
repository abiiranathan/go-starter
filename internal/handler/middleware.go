package handler

import (
	"net/http"

	"github.com/abiiranathan/rex"
	"github.com/abiiranathan/rex/middleware/logger"
	"github.com/abiiranathan/rex/middleware/recovery"
	"github.com/go-chi/chi/v5/middleware"
)

func AttachMiddleware(r *rex.Router) {
	r.Use(recovery.New(false, nil))
	r.Use(r.WrapMiddleware(middleware.RequestID))
	r.Use(r.WrapMiddleware(middleware.RealIP))

	cfg := logger.DefaultConfig
	cfg.Callback = func(r *http.Request, args ...any) []any {
		newArgs := []any{"Request ID", middleware.GetReqID(r.Context())}
		// ADD USER ID HERE
		return append(newArgs, args...)
	}

	r.Use(logger.New(cfg))

	r.Use(r.WrapMiddleware(middleware.Heartbeat("/ping")))
}
