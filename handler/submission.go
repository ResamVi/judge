package handler

import (
	"fmt"
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
	code, err := k.db.GetCode(c.Request().Context(), db.GetCodeParams{
		UserID:     user.ID,
		ExerciseID: task,
	})
	if err != nil || !code.Valid {
		slog.Error("could not find exercise for user", "task", task, "user", user.Username)
		return c.String(http.StatusNotFound, err.Error())
	}

	fmt.Println(code.String)

	data := k.page
	data.Body = fmt.Sprintf("%s%s%s", "<pre><code>",code.String,"</code></pre>")

	return c.Render(http.StatusOK, "index", data)
}
