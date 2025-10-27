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
	"time"
)

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

	return c.HTML(http.StatusOK, `<div><meta http-equiv="refresh" content="01; url=/">Erfolgreich angemeldet</div>`)
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

	err = k.db.CreateUser(c.Request().Context(), db.CreateUserParams{
		Username: username,
		Password: string(encrypted),
		Token:    uuid.New().String(),
	})
	if err != nil {
		slog.Error("Register: could not create user", "error", err)
		return c.NoContent(http.StatusBadRequest)
	}

	return c.HTML(http.StatusOK, `<span style="color:green;">Erfolgreich registriert.</span><br><a href="/login">Jetzt Anmelden</a>`)
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
