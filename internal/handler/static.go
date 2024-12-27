package handler

import (
	"time"

	"github.com/abiiranathan/go-starter/assets"
	"github.com/abiiranathan/rex"
)

func StaticRoutes(r *rex.Router) {
	maxAge := time.Hour * 24 * 30 // 30 days
	r.StaticFS("/static", rex.CreateFileSystem(assets.Static, "static"), int(maxAge.Seconds()))
}
