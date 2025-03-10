package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"text/template"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/yuin/goldmark"
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
		b, err := os.ReadFile("tasks/task01/README.md")
		if err != nil {
			panic(err)
		}

		var buf bytes.Buffer
		if err := goldmark.Convert(b, &buf); err != nil {
			panic(err)
		}

		return c.Render(http.StatusOK, "hello", buf.String())
	})
	e.Logger.Fatal(e.Start(":8080"))
}

type Template struct {
    templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
