package main

import (
	"io"
	"net/http"
	"text/template"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)
func main() {
	t := &Template{
		templates: template.Must(template.ParseGlob("www/index.html")),
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Renderer = t

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/render", func(c echo.Context) error {
		return c.Render(http.StatusOK, "hello", "world")
	})
	e.Logger.Fatal(e.Start(":8080"))
}

type Template struct {
    templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
