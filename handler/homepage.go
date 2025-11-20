package handler

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ResamVi/judge/grading"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

func (k Handler) Homepage(c echo.Context) error {
	// Contents of homepage comes from README.md file
	exerciseMD, err := os.ReadFile("www/README.md")
	if err != nil {
		slog.Error("os.ReadFile: " + err.Error())
		return err
	}

	// Convert local markdown files to HTML
	var exerciseHTML bytes.Buffer
	if err := md.Convert(exerciseMD, &exerciseHTML); err != nil {
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
		result = strings.ReplaceAll(exerciseHTML.String(), "{{Status}}", statusHTML)

		// Replace {{Token}} in webpage with user's token
		user, err := k.db.GetUser(c.Request().Context(), cookie.Value)
		if err != nil {
			slog.Error("failed to get user", "error", err, "username", cookie.Value)
			return err
		}
		result = strings.ReplaceAll(result, "{{Token}}", user.Token)
	} else {
		result = strings.ReplaceAll(exerciseHTML.String(), "{{Status}}", "")
		result = strings.ReplaceAll(result, "{{Token}}", "<Ein Token wird hier sichtbar sein sobald du dich eingeloggt hast>")
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

	usersHTML := ""
	for _, user := range users {
		usersHTML += fmt.Sprintf("<th style=\"text-align:center\">%s</th>", user.Username)
	}

	exercises, err := k.db.GetExercises(ctx)
	if err != nil {
		return "", fmt.Errorf("getting exercises: %w", err)
	}

	// TODO: This would be prettier using template language (low)
	exercisesHTML := ""
	for _, exercise := range exercises {
		exercisesHTML += "<tr>"
		exercisesHTML += fmt.Sprintf(`<td><a href="%s">`+exercise.Title+`</a></th>`, "/exercises/"+exercise.ID)

		status, err := k.db.GetStatus(ctx, exercise.ID)
		if err != nil {
			return "", fmt.Errorf("getting solvers: %w", err)
		}

		for _, row := range status {
			switch grading.Grade(row.Solved) {
			case grading.NotAttempted:
				exercisesHTML += fmt.Sprintf(`<td style="text-align:center"><a href="/submission/%s/%d"> </a></td>`, exercise.ID, row.UserID)
			case grading.Attempted:
				exercisesHTML += fmt.Sprintf(`<td style="text-align:center"><a href="/submission/%s/%d">❌</a></td>`, exercise.ID, row.UserID)
			case grading.Solved:
				exercisesHTML += fmt.Sprintf(`<td style="text-align:center"><a href="/submission/%s/%d">✔️ (%d)</a></td>`, exercise.ID, row.UserID, row.Attempts.Int32)
			}
		}
		exercisesHTML += "</tr>"
	}

	return fmt.Sprintf(`
	<h1>Bisheriger Stand</h1>

	<table style="width:100%%">
	<tr>
		<td></td>
		%s
	</tr>
	%s
	</table>`, usersHTML, exercisesHTML), nil
}
