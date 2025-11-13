package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

func (k Handler) Token(c echo.Context) error {
	token := "&lt;Ein Token wird hier sichtbar sein sobald du dich eingeloggt hast&gt;"

	tpl := "<h2>Dein Token</h2>"
	tpl += "<p>Bevor wir anfangen muss einmal der Judge wissen wer du bist.</p>"
	tpl += "<p>Hierzu kopierst du diesen Token in die Befehlszeile:</p>"
	tpl += "<pre><code>%s</code></pre>"
	tpl += "<p>Danach lassen sich Übungen herunterladen und hochladen.</p>"

	cookie, err := c.Cookie("username")
	if err != nil {
		data := k.page
		data.Body = fmt.Sprintf(tpl, token)
		return c.Render(http.StatusOK, "index", data)
	}

	if user, err := k.db.GetUser(c.Request().Context(), cookie.Value); err == nil {
		token = user.Token
	}

	data := k.page
	data.Body = fmt.Sprintf(tpl, token)
	return c.Render(http.StatusOK, "index", data)
}

func (k Handler) ValidateToken(c echo.Context) error {
	token := c.Request().Header.Get("token")

	_, err := k.db.GetUserFromToken(c.Request().Context(), token)
	if err != nil {
		slog.Warn("ValidateToken: user not found with token", "token", token)
		return c.NoContent(http.StatusNotFound)
	}

	return c.NoContent(http.StatusOK)
}
