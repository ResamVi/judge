package handler

import (
	"database/sql"
	"errors"
	"github.com/ResamVi/judge/db"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	"regexp"
	"time"
)

func (k Handler) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := k.db.GetUser(c.Request().Context(), username)
	if err != nil {
		slog.Warn("Login: user not found", "username", username, "error", err)
		return c.NoContent(http.StatusBadRequest)
	}

	if !user.Approved && k.env == "production" {
		slog.Warn("Login: User not approved", "username", username)
		return c.HTML(http.StatusOK, `<span style="color:red;">Benutzer noch nicht genehmigt.</span><a href="/login"></a>`)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		slog.Warn("Login: password does not match", "username", username, "error", err)
		return c.HTML(http.StatusOK, `<span style="color:red;">Anmeldung fehlgeschlagen.</span><a href="/login"></a>`)
	}

	cookie := &http.Cookie{
		Name:     "username",
		Value:    username,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   true,
		HttpOnly: false,
	}
	c.SetCookie(cookie)

	return c.HTML(http.StatusOK, `<div><meta http-equiv="refresh" content="01; url=/">Erfolgreich angemeldet</div>`)
}

func (k Handler) Register(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	confirm := c.FormValue("confirm")

	_, err := k.db.GetUser(c.Request().Context(), username)
	if password != confirm {
		slog.Warn("Register: password and confirmation do not match")
		return c.HTML(http.StatusOK, `<span style="color:red;">Registrierung fehlgeschlagen: Passwörter stimmen nicht überein</span><br><a href="/register">Zurück</a>`)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		slog.Warn("Register: username exists already", "username", username, "error", err)
		return c.HTML(http.StatusOK, `<span style="color:red;">Registrierung fehlgeschlagen: Benutzername existiert schon</span><br><a href="/register">Zurück</a>`)
	}

	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Register: could not generate password", "error", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	if len(username) > 12 || !regexp.MustCompile("^[a-zA-Z0-9]+$").MatchString(username) {
		slog.Error("Register: username disallowed", "username", username)
		return c.HTML(http.StatusOK, `<span style="color:red;">Registrierung fehlgeschlagen: Benutzername nicht länger als 12 Zeichen und nur Buchstaben und Zahlen</span><br><a href="/register">Zurück</a>`)
	}

	approved := false
	if k.env == "development" {
		approved = true
	}

	err = k.db.CreateUser(c.Request().Context(), db.CreateUserParams{
		Username: username,
		Password: string(encrypted),
		Token:    uuid.New().String(),
		Approved: approved,
	})
	if err != nil {
		slog.Error("Register: could not create user", "error", err)
		return c.NoContent(http.StatusBadRequest)
	}

	return c.HTML(http.StatusOK, `<span style="color:green;">Erfolgreich registriert.</span><br><a href="/login">Jetzt Anmelden</a>`)
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

func (k Handler) ValidateUsername(c echo.Context) error {
	username := c.FormValue("username")

	_, err := k.db.GetUser(c.Request().Context(), username)

	if !errors.Is(err, sql.ErrNoRows) {
		slog.Warn("ValidateUsername: username already taken", "username", username)

		return c.HTML(http.StatusOK, `
			<div style='color:red;'>Der Benutzername ist bereits vergeben.</div>
			<input name="username" hx-post="/validate/name" hx-target="#username-form" hx-indicator="#ind" value=`+username+`>
		`)
	}

	if len(username) > 12 {
		slog.Warn("ValidateUsername: username too long", "username", username)

		return c.HTML(http.StatusOK, `
			<div style='color:red;'>Der Benutzername darf höchstens 12 Zeichen haben..</div>
			<input name="username" hx-post="/validate/name" hx-target="#username-form" hx-indicator="#ind" value=`+username+`>
		`)

	}

	if !regexp.MustCompile("^[a-zA-Z0-9]+$").MatchString(username) {
		slog.Warn("ValidateUsername: username has invalid characters", "username", username)

		return c.HTML(http.StatusOK, `
			<div style='color:red;'>Der Benutzername nur Buchstaben und Zahlen haben.</div>
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
		slog.Warn("ValidatePassword: length too short") // field with length

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
		slog.Warn("ValidateConfirmation: passwords do not match")

		return c.HTML(http.StatusOK, `
			<div style="color:red">Passwörter müssen übereinstimmen.</div>
			<input name="confirm" type="password" hx-post="/validate/confirm" hx-target="#confirm-form" value="`+confirm+`">
		`)
	}

	return c.HTML(http.StatusOK, `
		<input name="confirm" type="password" hx-post="/validate/confirm" hx-target="#confirm-form" value="`+confirm+`">
	`)
}
