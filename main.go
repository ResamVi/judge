package main

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/plouc/textree"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
	"golang.org/x/crypto/bcrypt"

	"github.com/ResamVi/judge/db"
	"github.com/ResamVi/judge/handler"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	base = template.Must(template.ParseGlob("www/index.html"))
	md = goldmark.New(goldmark.WithRendererOptions(html.WithUnsafe()))
)
func main() {
	url := "postgres:postgres@localhost:5432/mydb?sslmode=disable" // TODO

	m, err := migrate.New("file://migrations", "pgx://"+url)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgres://"+url)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	encrypted, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost) // TODO: lol
	if err != nil {
		panic(err)
	}

	err = queries.UpsertUser(ctx, db.UpsertUserParams{
		Username: "admin",
		Password: string(encrypted),
		Approved: true,
	})
	if err != nil {
		panic(err)
	}

	t := &Template{
		templates: template.Must(template.ParseGlob("www/index.html")),
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Renderer = t

	h := handler.New(e, queries)

	e.GET("/", Homepage)
	e.GET("/login", LoginView)
	e.GET("/register", RegisterView)

	e.GET("/tasks/:task", TaskHandler)
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

func Homepage(c echo.Context) error {
	// Convert local markdown files to HTML
	taskMD, err := os.ReadFile("www/README.md")
	if err != nil {
		return err
	}

	var taskHTML bytes.Buffer
	if err := md.Convert(taskMD, &taskHTML); err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := base.ExecuteTemplate(&buf, "index", taskHTML.String()); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, buf.String())
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

func TaskHandler(c echo.Context) error {
	// Convert local markdown files to HTML
	taskMD, err := os.ReadFile("tasks/" + c.Param("task") + "/README.md")
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	var taskHTML bytes.Buffer
	if err := md.Convert(taskMD, &taskHTML); err != nil {
		return err
	}

	// Put that converted markdown into the webpage template for display
	var buf bytes.Buffer
	if err := base.ExecuteTemplate(&buf, "index", taskHTML.String()); err != nil {
		return err
	}

	// Replace occurrences of {{Code}} in the webpage with a custom file viewer
	htm := fmt.Sprintf(`
		<div class="row" style="margin-top: 3em; margin-bottom:3em">
			<div class="col-3">
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

		strippedPath := strings.TrimPrefix(path, "tasks/" + c.Param("task") + "/code/")

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

	err := filepath.Walk("tasks/" + c.Param("task") + "/code", walker)
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
