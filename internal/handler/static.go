package handler

import (
	"github.com/abiiranathan/go-starter/assets"
	"github.com/abiiranathan/rex"
)

func StaticRoutes(r *rex.Router) {
	r.StaticFS("/static", rex.CreateFileSystem(assets.Static, "static"), 0)
}
