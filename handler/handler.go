package handler

import (
	"bytes"
	"fmt"
	"github.com/ResamVi/judge/db"
	"github.com/labstack/echo/v4"
	"github.com/plouc/textree"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

var (
	md = goldmark.New(goldmark.WithRendererOptions(html.WithUnsafe()))
)

// Page contains all the data displayed on the current page
type Page struct {
	// List of exercises
	Exercises []MenuItem

	// The main content (derived from markdown files in tasks/)
	Body string
}

type MenuItem struct {
	Name string
	Link string
}

type Handler struct {
	db   *db.Queries
	page Page

	env string
}

func New(queries *db.Queries, env string) (*Handler, error) {
	exercises, err := getExercises()
	if err != nil {
		return nil, fmt.Errorf("getting exercises: %w", err)
	}

	return &Handler{
		db: queries,
		page: Page{
			Exercises: exercises,
		},
		env: env,
	}, nil
}

func (k Handler) TaskHandler(c echo.Context) error {
	// Convert local markdown files to HTML
	taskMD, err := os.ReadFile("tasks/" + c.Param("task") + "/README.md")
	if err != nil {
		slog.Error("os.ReadFile: " + err.Error())
		return c.NoContent(http.StatusNotFound)
	}
	var taskHTML bytes.Buffer
	if err := md.Convert(taskMD, &taskHTML); err != nil {
		slog.Error("md.Convert: " + err.Error())
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

	result := strings.ReplaceAll(taskHTML.String(), "{{Code}}", htm)

	data := k.page
	data.Body = result

	return c.Render(http.StatusOK, "index", data)
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

// Username shows the name if the user is logged and if not a link to login
func (k Handler) Username(c echo.Context) error {
	cookie, err := c.Cookie("username")
	if err != nil {
		slog.Info("user not logged in: " + err.Error())
		return c.HTML(http.StatusOK, `<li class="float-right"><a href="/login">Anmelden</a></li>`)
	}

	return c.HTML(http.StatusOK, `<li class="float-right">Eingeloggt als <strong>`+cookie.Value+`</strong></li>`)
}

func (k Handler) LoginView(c echo.Context) error {
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

	data := k.page
	data.Body = str

	return c.Render(http.StatusOK, "index", data)
}

func (k Handler) RegisterView(c echo.Context) error {
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

	data := k.page
	data.Body = str

	return c.Render(http.StatusOK, "index", data)
}

func getExercises() ([]MenuItem, error) {
	entries, err := os.ReadDir("tasks")
	if err != nil {
		return nil, fmt.Errorf("read dir of tasks: %w", err)
	}

	var exercises []MenuItem
	for _, e := range entries {
		file, err := os.ReadFile("tasks/" + e.Name() + "/README.md")
		if err != nil {
			return nil, fmt.Errorf("read README.md of folder: %w", err)
		}
		title, _, found := bytes.Cut(file, []byte("\n"))
		if !found {
			return nil, fmt.Errorf("missing title in README.md of " + e.Name())
		}
		title = bytes.TrimPrefix(title, []byte("# "))

		name, _, _ := strings.Cut(e.Name(), "-")

		exercises = append(exercises, MenuItem{
			Name: fmt.Sprintf("Aufgabe %s: %s", name, string(title)),
			Link: "/tasks/" + e.Name(),
		})
	}

	return exercises, nil
}
