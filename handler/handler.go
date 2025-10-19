package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"bytes"
	"github.com/ResamVi/judge/db"
	"github.com/labstack/echo/v4"
	"github.com/plouc/textree"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"os"
	"strings"
)

var (
	md   = goldmark.New(goldmark.WithRendererOptions(html.WithUnsafe()))
	base = template.Must(template.ParseGlob("www/index.html"))
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

func (k Handler) Homepage(c echo.Context) error {
	// Contents of homepage comes from README.md file
	taskMD, err := os.ReadFile("www/README.md")
	if err != nil {
		slog.Error("os.ReadFile: " + err.Error())
		return err
	}

	// Convert local markdown files to HTML
	var taskHTML bytes.Buffer
	if err := md.Convert(taskMD, &taskHTML); err != nil {
		slog.Error("md.Convert: " + err.Error())
		return err
	}

	data := k.page
	data.Body = taskHTML.String()

	return c.Render(http.StatusOK, "index", data)
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

func (k Handler) ValidateUsername(c echo.Context) error {
	username := c.FormValue("username")

	_, err := k.db.GetUser(c.Request().Context(), username)

	if !errors.Is(err, sql.ErrNoRows) {
		slog.Error("ValidateUsername: username already taken") // TODO: show which username is taken

		return c.HTML(http.StatusOK, `
			<div style='color:red;'>Der Benutzername ist bereits vergeben.</div>
			<input name="username" hx-post="/validate/name" hx-target="#username-form" hx-indicator="#ind" value=`+username+`>
		`)
	}

	return c.HTML(http.StatusOK, `
		<input name="username" hx-post="/validate/name" hx-target="#username-form" hx-indicator="#ind" value=`+username+`>
	`)
}

func (k Handler) ValidatePassword(c echo.Context) error {
	password := c.FormValue("password")

	if len(password) < 8 {
		slog.Error("ValidatePassword: length too short") // field with length

		return c.HTML(http.StatusOK, `
			<div style="color:red">Password muss mindestens 8 Zeichen haben.</div>
			<input name="password" type="password" hx-post="/validate/password" hx-target="#password-form" value="`+password+`">
		`)
	}

	return c.HTML(http.StatusOK, `
		<input name="password" type="password" hx-post="/validate/password" hx-target="#password-form" value="`+password+`">
	`)
}

func (k Handler) ValidateConfirmation(c echo.Context) error {
	confirm := c.FormValue("confirm")
	password := c.FormValue("password")

	if confirm != password {
		slog.Error("ValidateConfirmation: passwords do not match")

		return c.HTML(http.StatusOK, `
			<div style="color:red">Passwörter müssen übereinstimmen.</div>
			<input name="confirm" type="password" hx-post="/validate/confirm" hx-target="#confirm-form" value="`+confirm+`">
		`)
	}

	return c.HTML(http.StatusOK, `
		<input name="confirm" type="password" hx-post="/validate/confirm" hx-target="#confirm-form" value="`+confirm+`">
	`)
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

func (k Handler) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := k.db.GetUser(c.Request().Context(), username)
	if err != nil {
		slog.Error("Login: user not found") // TODO: field of username
		return c.NoContent(http.StatusBadRequest)
	}

	if !user.Approved && k.env == "production" {
		slog.Error("Login: User not approved", "username", username)
		return c.HTML(http.StatusOK, `<span style="color:red;">Benutzer noch nicht genehmigt.</span><a href="/login"></a>`)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		slog.Error("Login: password does not match") // TODO: error of username
		return c.HTML(http.StatusOK, `<span style="color:red;">Anmeldung fehlgeschlagen.</span><a href="/login"></a>`)
	}

	cookie := &http.Cookie{
		Name:     "username",
		Value:    username,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
	}
	c.SetCookie(cookie)

	return c.HTML(http.StatusOK, `<div><meta http-equiv="refresh" content="1; url=/">Erfolgreich angemeldet</div>`)
}

func (k Handler) Register(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	confirm := c.FormValue("confirm")

	_, err := k.db.GetUser(c.Request().Context(), username)
	if password != confirm {
		slog.Warn("Register: password and confirmation do not match")
		return c.HTML(http.StatusOK, `<span style="color:red;">Registrierung fehlgeschlagen</span><br><a href="/register">Zurück</a>`)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		slog.Warn("Register: username exists already", "username", username, "error", err)
		return c.HTML(http.StatusOK, `<span style="color:red;">Registrierung fehlgeschlagen</span><br><a href="/register">Zurück</a>`)
	}

	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Register: could not generate password", "error", err)
		return c.NoContent(http.StatusBadRequest)
	}

	_, err = k.db.CreateUser(c.Request().Context(), db.CreateUserParams{
		Username: username,
		Password: string(encrypted),
	})
	if err != nil {
		slog.Error("Register: could not create user", "error", err)
		return c.NoContent(http.StatusBadRequest)
	}

	return c.HTML(http.StatusOK, `<span style="color:green;">Erfolgreich registriert.</span><br><a href="/login">Jetzt Anmelden</a>`)
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

		exercises = append(exercises, MenuItem{
			Name: fmt.Sprintf("Aufgabe %s: %s", e.Name(), string(title)),
			Link: "/tasks/" + e.Name(),
		})
	}

	return exercises, nil
}
