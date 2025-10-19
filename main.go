package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ResamVi/judge/db"
	"github.com/ResamVi/judge/handler"
	"github.com/ResamVi/judge/migrate"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/plouc/textree"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

var (
	base = template.Must(template.ParseGlob("www/index.html"))
	md   = goldmark.New(goldmark.WithRendererOptions(html.WithUnsafe()))
)

func main() {
	url := "postgres:postgres@localhost:5432/mydb?sslmode=disable" // TODO

	err := migrate.DB(url)
	if err != nil {
		panic(err)
	}

	queries, err := db.Init(url)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("www/index.html")),
	}

	h := handler.New(queries)

	e.GET("/", h.Homepage)
	e.GET("/login", LoginView)
	e.GET("/register", RegisterView)

	e.GET("/tasks/:task", h.TaskHandler)
	e.GET("/tasks/:task/code", CodeHandler)

	e.POST("/login", h.Login)
	e.POST("/register", h.Register)
	e.POST("/username", h.Username)
	e.POST("/validate/name", h.ValidateUsername)
	e.POST("/validate/password", h.ValidatePassword)
	e.POST("/validate/confirm", h.ValidateConfirmation)

	e.Static("/tasks", "tasks")
	e.Static("/www", "www")

	e.Logger.Fatal(e.Start(":8080"))
}

func LoginView(c echo.Context) error {
	str := `
	<form style="margin-top: 2em" method="post" action="/login">
	<div class="container">
		<label for="uname"><b>Benutzername</b></label>
		<input type="text" placeholder="Benutzername" name="username" required>

		<label for="psw"><b>Passwort</b></label>
		<input type="password" placeholder="Passwort" name="password" required>

		<button type="submit" >Login</button>
		<!-- <label><input type="checkbox" checked="checked" name="remember"> Eingeloggt bleiben</label> TODO: -->
	</div>

	<div class="container" style="background-color:#f1f1f1">
		<span class="psw">Noch nicht <a href="/register">registriert</a>?</span>
	</div>
	</form> 
	`

	var buf bytes.Buffer
	if err := base.ExecuteTemplate(&buf, "index", str); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, buf.String())
}

func RegisterView(c echo.Context) error {
	str := `
	<form style="margin-top: 2em" hx-post="/register">
	<div hx-target="this">
		<label><b>Benutzername</b></label>
		<div id="username-form">
			<input name="username" hx-post="/user/name" hx-target="#username-form" hx-indicator="#ind">
		</div>

		<label><b>Passwort</b></label>
		<div id="password-form">
			<input name="password" type="password" hx-post="/user/password" hx-target="#password-form">
		</div>

		<label><b>Passwort bestätigen</b></label>
		<div id="confirm-form">
			<input name="confirm" type="password" hx-post="/user/confirm" hx-target="#confirm-form">
		</div>

		<button class="btn primary">Registrieren</button>
		<img id="ind" src="/assets/bars.svg" class="htmx-indicator"/>
	</div>
	</form>
	`

	var buf bytes.Buffer
	if err := base.ExecuteTemplate(&buf, "index", str); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, buf.String())
}

func CodeHandler(c echo.Context) error {
	var b []byte
	buf := bytes.NewBuffer(b)
	w := zip.NewWriter(buf)

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		strippedPath := strings.TrimPrefix(path, "tasks/"+c.Param("task")+"/code/")

		f, err := w.Create(strippedPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}

	err := filepath.Walk("tasks/"+c.Param("task")+"/code", walker)
	if err != nil {
		panic(err)
	}
	w.Close()

	return c.Blob(http.StatusOK, "application/zip", buf.Bytes())
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
			<input type="radio" name="tabs" id="tabone" checked="checked">
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
