package handlers

import (
	"net/http"

	"github.com/abiiranathan/rex"
	"github.com/abiiranathan/rex/middleware/logger"
	"github.com/abiiranathan/rex/middleware/recovery"
	"github.com/go-chi/chi/v5/middleware"
)

func AttachMiddleware(r *rex.Router) {
	const printStackTrace = false

	r.Use(recovery.New(printStackTrace))
	r.Use(r.WrapMiddleware(middleware.RequestID))
	r.Use(r.WrapMiddleware(middleware.RealIP))

	loggerCfg := logger.DefaultConfig
	loggerCfg.Callback = func(r *http.Request, args ...any) []any {
		newArgs := []any{"Request ID", middleware.GetReqID(r.Context())}
		// newArgs = append(newArgs, "user_id", 1)
		return append(newArgs, args...)
	}

	r.Use(logger.New(loggerCfg))
}
