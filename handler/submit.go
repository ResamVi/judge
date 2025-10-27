package handler

import (
	"encoding/base64"
	"github.com/labstack/echo/v4"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func (k Handler) Submit(c echo.Context) error {
	token := c.Request().Header.Get("token")
	exercise := c.Request().Header.Get("exercise")

	user, err := k.db.GetUserFromToken(c.Request().Context(), token)
	if err != nil {
		slog.Error("could not find user with token", "token", token, "err", err)
		return c.String(http.StatusNotFound, err.Error())
	}

	// todo: add exercise no and attempt count
	f, err := os.Create("submissions/" + time.Now().Format("2006-01-02T15-04") + "_" + user.Username + exercise + ".zip")
	if err != nil {
		slog.Error("could not create file in submission folder", "err", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	_, err = io.Copy(f, base64.NewDecoder(base64.StdEncoding, c.Request().Body))
	if err != nil {
		slog.Error("could not copy http body to file", "err", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, ":)")
}
