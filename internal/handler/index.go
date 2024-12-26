package handler

import (
	"github.com/abiiranathan/go-starter/templates"
	"github.com/abiiranathan/rex"
)

func IndexRoute(h *handler) {
	h.router.GET("/", func(c *rex.Context) error {
		return RenderPage(c, templates.Index())
	})

	h.router.GET("/offline", func(c *rex.Context) error {
		return RenderPage(c, templates.OfflinePage())
	})
}
