package handlers

import (
	"github.com/abiiranathan/rex"
)

func IndexRoute(h *handler) {
	h.router.GET("/", func(c *rex.Context) error {
		return c.Render("templates/index", rex.Map{})
	})
}
