package handler

import (
	"time"

	"github.com/abiiranathan/rex"
	"github.com/abiiranathan/rex/middleware/logger"
	"github.com/abiiranathan/rex/middleware/recovery"
	"github.com/go-chi/chi/v5/middleware"
)

func AttachMiddleware(r *rex.Router) {
	r.Use(recovery.New(false, nil))
	r.Use(r.WrapMiddleware(middleware.RequestID))
	r.Use(r.WrapMiddleware(middleware.RealIP))
	r.Use(logger.New(logger.DefaultConfig))

	r.Use(r.WrapMiddleware(middleware.Timeout(time.Second * 30)))
	r.Use(r.WrapMiddleware(middleware.Heartbeat("/ping")))
}
