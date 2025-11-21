package handler

import (
	"log/slog"
	"net/http"

	"fmt"
	"github.com/ResamVi/judge/db"
	"github.com/ResamVi/judge/grading"
	"github.com/labstack/echo/v4"
	"strconv"
)

func (k Handler) Submission(c echo.Context) error {
	cookie, err := c.Cookie("username")
	if err != nil {
		slog.Warn("must be logged in to see suggestions", "error", err)
		return fmt.Errorf("must be logged in to see suggestions	")
	}

	exercise := c.Param("exercise")
	viewer, err := k.db.GetUser(c.Request().Context(), cookie.Value)
	if err != nil {
		slog.Warn("unknown user", "username", cookie.Value, "error", err)
		return fmt.Errorf("must be logged in to see suggestions	")
	}
	viewerSubm, err := k.db.GetSubmission(c.Request().Context(), db.GetSubmissionParams{
		UserID:     viewer.ID,
		ExerciseID: exercise,
	})
	if err != nil {
		data := k.page
		data.Body += "Aufgabe muss gelöst werden bevor man andere Lösungen ansehen darf"
		return c.Render(http.StatusOK, "index", data)
	}
	if fmt.Sprintf("%d", viewer.ID) != c.Param("user") && viewerSubm.Solved != int32(grading.Solved) {
		data := k.page
		data.Body += "Aufgabe muss gelöst werden bevor man andere Lösungen ansehen darf"
		return c.Render(http.StatusOK, "index", data)
	}

	userId, err := strconv.Atoi(c.Param("user"))
	if err != nil {
		slog.Error("could not parse user id", "user param", c.Param("user"), "error", err)
		return fmt.Errorf("could not parse user id")
	}

	user, err := k.db.GetUserFromId(c.Request().Context(), int64(userId))
	if err != nil {
		slog.Error("could not find user", "userId", userId, "error", err)
		return c.String(http.StatusNotFound, err.Error())
	}

	subm, err := k.db.GetSubmission(c.Request().Context(), db.GetSubmissionParams{
		UserID:     user.ID,
		ExerciseID: exercise,
	})
	if err != nil {
		slog.Error("could not find submission for user", "exercise", exercise, "user", user.Username)
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
