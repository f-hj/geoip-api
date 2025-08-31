package main

import (
	"embed"
	"html/template"
	"io"

	"github.com/labstack/echo"
)

//go:embed templates
var templates embed.FS

type Template struct {
	templates *template.Template
}

func NewHTMLTemplates() (*Template, error) {
	templates, err := template.ParseFS(templates, "templates/*.html")
	if err != nil {
		return nil, err
	}
	return &Template{
		templates,
	}, nil
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
