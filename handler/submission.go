package handler

import (
	"log/slog"
	"net/http"

	"fmt"
	"github.com/ResamVi/judge/db"
	"github.com/labstack/echo/v4"
	"strconv"
)

func (k Handler) Submission(c echo.Context) error {
	userId, err := strconv.Atoi(c.Param("user"))
	if err != nil {
		return fmt.Errorf("could not parse user id")
	}

	user, err := k.db.GetUserFromId(c.Request().Context(), int64(userId))
	if err != nil {
		slog.Error("could not find user", "userId", userId)
		return c.String(http.StatusNotFound, err.Error())
	}

	task := c.Param("task")
	subm, err := k.db.GetSubmission(c.Request().Context(), db.GetSubmissionParams{
		UserID:     user.ID,
		ExerciseID: task,
	})
	if err != nil {
		slog.Error("could not find submission for user", "task", task, "user", user.Username)
		return c.HTML(http.StatusOK, `<div><meta http-equiv="refresh" content="03; url=/">Benutzer hat Aufgabe noch nicht gelöst</div>`)
	}

	data := k.page
	data.Body += "<h2>Output</h2>"
	data.Body += "<pre><code>" + subm.Output + "</code></pre>"
	data.Body += "<h2>Evaluation</h2>"
	data.Body += subm.Evaluation
	data.Body += "<h2>Code</h2>"
	data.Body += "<pre><code>" + subm.Code + "</code></pre>"

	return c.Render(http.StatusOK, "index", data)
}
