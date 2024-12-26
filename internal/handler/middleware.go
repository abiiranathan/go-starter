package handler

import (
	"time"

	"github.com/abiiranathan/rex"
	"github.com/go-chi/chi/v5/middleware"
)

func AttachMiddleware(r *rex.Router) {
	r.Use(r.WrapMiddleware(middleware.RequestID))
	r.Use(r.WrapMiddleware(middleware.RealIP))
	r.Use(r.WrapMiddleware(middleware.Heartbeat("/ping")))
	r.Use(r.WrapMiddleware(middleware.Timeout(time.Second * 30)))
	r.Use(r.WrapMiddleware(middleware.RedirectSlashes))
}
