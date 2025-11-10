package handler

import (
	"bytes"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

func (k Handler) Homepage(c echo.Context) error {
	// Contents of homepage comes from README.md file
	taskMD, err := os.ReadFile("www/README.md")
	if err != nil {
		slog.Error("os.ReadFile: " + err.Error())
		return err
	}

	// Convert local markdown files to HTML
	var taskHTML bytes.Buffer
	if err := md.Convert(taskMD, &taskHTML); err != nil {
		slog.Error("md.Convert: " + err.Error())
		return err
	}

	// Status should only appear when logged in
	var result string
	if cookie, err := c.Cookie("username"); err == nil {
		// Replace {{Status}} in webpage with table of user progress
		statusHTML, err := k.status(c.Request().Context())
		if err != nil {
			slog.Error("failed to generate status HTML", "error", err)
			return err
		}
		result = strings.ReplaceAll(taskHTML.String(), "{{Status}}", statusHTML)

		// Replace {{Token}} in webpage with user's token
		user, err := k.db.GetUser(c.Request().Context(), cookie.Value)
		if err != nil {
			slog.Error("failed to get user", "error", err, "username", cookie.Value)
			return err
		}
		result = strings.ReplaceAll(result, "{{Token}}", user.Token)
	} else {
		result = strings.ReplaceAll(taskHTML.String(), "{{Status}}", "")
		result = strings.ReplaceAll(result, "{{Token}}", "<my-token>")
	}

	data := k.page
	data.Body = result

	return c.Render(http.StatusOK, "index", data)
}

func (k Handler) status(ctx context.Context) (string, error) {
	users, err := k.db.GetUsers(ctx)
	if err != nil {
		return "", fmt.Errorf("getting users: %w", err)
	}

	exercises, err := k.db.GetExercises(ctx)
	if err != nil {
		return "", fmt.Errorf("getting exercises: %w", err)
	}

	usersHTML := ""
	for _, user := range users {
		usersHTML += fmt.Sprintf("<th>%s</th>", user.Username)
	}

	// TODO: Clickable links
	exercisesHTML := ""
	for _, exercise := range exercises {
		exercisesHTML += "<tr>"
		id, _, _ := strings.Cut(exercise.ID, "-")
		exercisesHTML += fmt.Sprintf("<th>Aufgabe %s: %s</td>", id, exercise.Title)

		solvers, err := k.db.GetSolvers(ctx, exercise.ID)
		if err != nil {
			return "", fmt.Errorf("getting solvers: %w", err)
		}

		for _, solver := range solvers {
			if solver.Solved {
				exercisesHTML += fmt.Sprintf("<td>✔️</td>")
			} else {
				exercisesHTML += fmt.Sprintf("<td>❌</td>")
			}
		}
		exercisesHTML += "</tr>"
	}

	return fmt.Sprintf(`
	<h1>Status</h1>

	<table>
	<tr>
		<td></td>
		%s
	</tr>
	%s
	</table>`, usersHTML, exercisesHTML), nil
}
