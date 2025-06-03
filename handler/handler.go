package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/ResamVi/judge/db"
	"github.com/labstack/echo"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
	"golang.org/x/crypto/bcrypt"
)

var(
	md = goldmark.New(goldmark.WithRendererOptions(html.WithUnsafe()))
	base = template.Must(template.ParseGlob("www/index.html"))
) 

type Handler struct {
	db *db.Queries
}

func New(e *echo.Echo, queries *db.Queries) *Handler {
	return &Handler{
		db: queries,
	}
}

func (k Handler) Username(c echo.Context) error {
	cookie, err := c.Cookie("username")
	if err != nil {
		return c.HTML(http.StatusOK, `<li class="float-right"><a href="/login">Anmelden</a></li>`)
	}

	return c.HTML(http.StatusOK, `<li class="float-right">Eingeloggt als <strong>` + cookie.Value + `</strong></li>`)
}

func (k Handler) ValidateUsername(c echo.Context) error {
	username := c.FormValue("username")

	_, err := k.db.GetUser(c.Request().Context(), username)

	if !errors.Is(err, sql.ErrNoRows) {
		return c.HTML(http.StatusOK, `
			<div style='color:red;'>Der Benutzername ist bereits vergeben.</div>
			<input name="username" hx-post="/validate/name" hx-target="#username-form" hx-indicator="#ind" value=` + username + `>
		`)
	}

	return c.HTML(http.StatusOK, `
		<input name="username" hx-post="/validate/name" hx-target="#username-form" hx-indicator="#ind" value=` + username + `>
	`)
}

func (k Handler) ValidatePassword(c echo.Context) error {
	password := c.FormValue("password")

	if len(password) < 8 {
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
		return c.HTML(http.StatusOK, `
			<div style="color:red">Passwörter müssen übereinstimmen.</div>
			<input name="confirm" type="password" hx-post="/validate/confirm" hx-target="#confirm-form" value="`+confirm+`">
		`)
	}

	return c.HTML(http.StatusOK, `
		<input name="confirm" type="password" hx-post="/validate/confirm" hx-target="#confirm-form" value="`+confirm+`">
	`)
}


func (k Handler) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	fmt.Println(username, password)

	user, err := k.db.GetUser(c.Request().Context(), username)
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println("login failed for '" + user.Username + "' " + err.Error())
		return c.HTML(http.StatusOK, `<span style="color:red;">Anmeldung fehlgeschlagen.</span><a href="/login"></a>`)
	}

	cookie := &http.Cookie{
		Name:        "username",
		Value:       username,
		Path:        "/",
		Expires:     time.Now().Add(24 * time.Hour),
		Secure:      true,
		HttpOnly:    true,
	}
	c.SetCookie(cookie)

	return c.HTML(http.StatusOK, `<div><meta http-equiv="refresh" content="1; url=/">Erfolgreich angemeldet</div>`)
}


func (k Handler) Register(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	confirm := c.FormValue("confirm")

	_, err := k.db.GetUser(c.Request().Context(), username)
	if password != confirm || !errors.Is(err, sql.ErrNoRows) {
		return c.HTML(http.StatusOK, `<span style="color:red;">Registrierung fehlgeschlagen</span><br><a href="/register">Zurück</a>`)
	}

	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}

	_, err = k.db.CreateUser(c.Request().Context(), db.CreateUserParams{
		Username: username,
		Password: string(encrypted),
	})
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}

	return c.HTML(http.StatusOK, `<span style="color:green;">Erfolgreich registriert.</span><br><a href="/login">Jetzt Anmelden</a>`)
}

