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

var (
	base = template.Must(template.ParseGlob("www/index.html"))
	md = goldmark.New(goldmark.WithRendererOptions(html.WithUnsafe()))
)

func main() {
	t := &Template{
		templates: template.Must(template.ParseGlob("www/index.html")),
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Renderer = t

	e.GET("/tasks/:task", taskHandler)
	e.GET("/code/:task", codeHandler)
	e.Static("/", "tasks")

	e.Logger.Fatal(e.Start(":8080"))
}

func codeHandler(c echo.Context) error {
	return nil
}

func taskHandler(c echo.Context) error {
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
		<div style="text-align:center"><button type="submit">Download</button></div>
		<div>%s</div>
		</div>

		<div class="col">

		<div class="tabs">
		%s
		</div>

		</div>
		</div>`, treeView(c.Param("task")), codeView(c.Param("task")))

	result := strings.ReplaceAll(buf.String(), "{{Code}}", htm)

	return c.HTML(http.StatusOK, result)
}

func treeView(name string) string {
	tree, err := textree.TreeFromDir("./tasks/" + name + "/code")
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
	entries, err := os.ReadDir("./tasks/" + name + "/code")
	if err != nil {
		panic(err)
	}

	code := ""
	for _, e := range entries {
		content, err := os.ReadFile("./tasks/" + name + "/code/" + e.Name())
		if err != nil {
			panic(err)
		}

		code += fmt.Sprintf(`
			<input type="radio" name="tabs" id="tabone">
			<label for="tabone">%s</label>
			<div class="tab"><pre><code>%s</code></pre></div>`, e.Name(), string(content))
	}

	return code
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
