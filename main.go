package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/plouc/textree"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

func main() {
	base := template.Must(template.ParseGlob("www/index.html"))

	t := &Template{
		templates: template.Must(template.ParseGlob("www/index.html")),
	}

	md := goldmark.New(goldmark.WithRendererOptions(html.WithUnsafe()))

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Renderer = t

	e.GET("/tasks/:task", func(c echo.Context) error {

		// Convert local markdown files to HTML
		taskMD, err := os.ReadFile("tasks/" + c.Param("task") + "/README.md")
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}

		var taskHTML bytes.Buffer
		if err := md.Convert(taskMD, &taskHTML); err != nil {
			return err
		}

		var buf bytes.Buffer
		if err := base.ExecuteTemplate(&buf, "index", taskHTML.String()); err != nil {
			return err
		}

		htm := fmt.Sprintf(`
<div class="row" style="margin-top: 3em; margin-bottom:3em">
  <div class="col-3">
    %s

  <button type="submit">Download der Aufgabe</button>
  </div>

  <div class="col">

	<div class="tabs">
		%s
	</div>

  </div>
</div>
		`, treeView(c.Param("task")), codeView(c.Param("task")))

		result := strings.ReplaceAll(buf.String(), "{{Code}}", htm)

		return c.HTML(http.StatusOK, result)
	})
	e.Static("/", "tasks")

	e.Logger.Fatal(e.Start(":8080"))
}

func treeView(name string) string {
	tree, err := textree.TreeFromDir("./code/" + name)
	if err != nil {
		panic(err)
	}

	var treebuf bytes.Buffer
	tree.Render(&treebuf, textree.NewRenderOptions())

	result := treebuf.String()
	result = strings.TrimSpace(result)
	result = strings.ReplaceAll(result, "\n", "<br />")

	return result
}

func codeView(name string) string {
	return `
	  <input type="radio" name="tabs" id="tabone" checked="checked">
	  <label for="tabone">Tab One</label>
	  <div class="tab"> some code </div>

	  <input type="radio" name="tabs" id="tabtwo">
	  <label for="tabtwo">Tab Two</label>
	  <div class="tab"> some code </div>
	  
	  <input type="radio" name="tabs" id="tabthree">
	  <label for="tabthree">Tab Three</label>
	  <div class="tab"> some code </div>`

}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
