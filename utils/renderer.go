package utils

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Templates struct {
	templates *template.Template
}

func NewTemplates() *Templates {
	t := template.New("")
	template.Must(t.ParseGlob("views/*.html"))
	template.Must(t.ParseGlob("views/components/*.html"))
	return &Templates{
		templates: t,
	}
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
