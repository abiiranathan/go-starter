package handlers

import (
	"time"

	"github.com/abiiranathan/go-starter/assets"
	"github.com/abiiranathan/rex"
)

func StaticRoutes(r *rex.Router) {
	// 30 days (0 is no cache)
	maxAge := int((time.Hour * 24 * 30).Seconds())

	r.StaticFS("/static", rex.CreateFileSystem(assets.Static, "static"), maxAge)
}
