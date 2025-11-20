package handler

import (
	"github.com/labstack/echo/v4"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func (k Handler) Editor(c echo.Context) error {
	username := c.Request().Header.Get("username")
	exercise := c.Request().Header.Get("exercise")

	user, err := k.db.GetUser(c.Request().Context(), username)
	if err != nil {
		slog.Error("could not find user with token", "username", username, "err", err)
		return c.String(http.StatusNotFound, err.Error())
	}
	// TODO: should check if approved

	code, err := io.ReadAll(c.Request().Body)
	if err != nil {
		slog.Error("bad body", "error", err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	destDir := filepath.Join("submissions", time.Now().Format("2006-01-02T15-04")+"_"+user.Username)

	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		slog.Error("cannot create dest dir", "error", err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = os.WriteFile(filepath.Join(destDir, "main.go"), code, os.ModePerm)
	if err != nil {
		slog.Error("cannot create main.go", "error", err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	// == Build & run code ==
	go k.executeSubmission(user, destDir, exercise, string(code))

	slog.Info("submission received (editor)", "user", user.Username, "exercise", exercise)
	return c.NoContent(http.StatusOK)

}
