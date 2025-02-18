package views

import (
	"embed"
	"html/template"

	"github.com/abiiranathan/rex"
)

//go:embed templates
var tmpl embed.FS

// Return embed templates inside views/templates as a parsed
// tree or error.
// Parse template functions as a funcMap.
func Templates(funcMap template.FuncMap) (*template.Template, error) {
	return rex.ParseTemplatesFS(tmpl, "templates", funcMap)
}
