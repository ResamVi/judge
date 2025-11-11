package handler

import (
	"log/slog"
	"net/http"

	"github.com/ResamVi/judge/db"
	"github.com/labstack/echo/v4"
)

func (k Handler) Submission(c echo.Context) error {
	cookie, err := c.Cookie("username")
	if err != nil {
		slog.Error("could not cookie", "error", err)
		return c.String(http.StatusNotFound, err.Error())
	}

	user, err := k.db.GetUser(c.Request().Context(), cookie.Value)
	if err != nil {
		slog.Error("could not find user", "cookie", cookie.Value)
		return c.String(http.StatusNotFound, err.Error())
	}

	task := c.Param("task")
	subm, err := k.db.GetSubmission(c.Request().Context(), db.GetSubmissionParams{
		UserID:     user.ID,
		ExerciseID: task,
	})
	if err != nil || !subm.Code.Valid {
		slog.Error("could not find submission for user", "task", task, "user", user.Username)
		return c.String(http.StatusNotFound, err.Error())
	}

	data := k.page
	data.Body = "<h2>Code</h2><pre><code>" + subm.Code.String + "</code></pre><h2>Output</h2><pre><code>" + subm.Output.String + "</code></pre>"

	return c.Render(http.StatusOK, "index", data)
}
