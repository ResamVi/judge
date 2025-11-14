package handler

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/ResamVi/judge/db"
	"github.com/labstack/echo/v4"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func (k Handler) Submit(c echo.Context) error {
	// == Unwrap user inputs ==
	token := c.Request().Header.Get("token")
	exercise := c.Request().Header.Get("exercise")

	user, err := k.db.GetUserFromToken(c.Request().Context(), token)
	if err != nil {
		slog.Error("could not find user with token", "token", token, "err", err)
		return c.String(http.StatusNotFound, err.Error())
	}

	data, err := io.ReadAll(base64.NewDecoder(base64.StdEncoding, c.Request().Body))
	if err != nil {
		slog.Error("bad base64", "error", err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	readerAt := bytes.NewReader(data)
	zr, err := zip.NewReader(readerAt, int64(len(data)))
	if err != nil {
		slog.Error("bad zip data", "error", err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	// == Store the code in database ==
	code := ""
	for _, f := range zr.File {
		code += fmt.Sprintf("\n// === File: %s ===\n", f.Name)

		rc, err := f.Open()
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		// Copy file contents to stdout
		b, err := io.ReadAll(rc)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		code += string(b)
		rc.Close()
	}

	// == Store code locally for execution ==
	destDir := filepath.Join("submissions", time.Now().Format("2006-01-02T15-04"))

	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		slog.Error("cannot create dest dir", "error", err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := unzipBytes(data, destDir); err != nil {
		slog.Error("unzip failed", "error", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// == Build & run code ==
	// IMPORTANT TODO: timeout for command
	runCmd := exec.Command("go", "run", destDir+"/main.go")

	var buildStderr, output bytes.Buffer
	runCmd.Stderr = &buildStderr
	runCmd.Stdout = &output

	if err := runCmd.Run(); err != nil {
		fmt.Printf("Build failed:\n%s\n", buildStderr.String())
		output = buildStderr
	}

	evaluation, solved := gradeSubmission(exercise, code, output.String())

	err = k.db.CreateSubmission(c.Request().Context(), db.CreateSubmissionParams{
		UserID:     user.ID,
		ExerciseID: exercise,
		Code:       code,
		Output:     output.String(),
		Evaluation: evaluation,
		Solved:     int32(solved),
	})
	if err != nil {
		slog.Error("failed to create submission", "userId", user.ID, "exercise", exercise, "error", err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	// TODO: Only keep 20 most recent submissions

	return c.NoContent(http.StatusOK)
}

// TODO: move to grading package
type Criteria struct {
	Description string
	Valid       func(code, output string) (comment string, failed bool)
}

type Exercise struct {
	Criteria []Criteria
}

var Grading = map[string]Exercise{
	"01-compiler": {
		Criteria: []Criteria{
			{
				Description: "Beispielkriterium",
				Valid: func(code, output string) (string, bool) {
					return "❌, Ist noch nicht am funktionieren", false
					// return "✅ Ist am funktionieren", true
				},
			},
		},
	},
}

func gradeSubmission(exercise string, code string, output string) (string, int) {
	criteria, ok := Grading[exercise]
	if !ok {
		return "Unbekannt: " + exercise, db.NotAttempted
	}

	evaluation := ""
	solved := db.Solved

	for _, fn := range criteria.Criteria {
		comment, valid := fn.Valid(code, output)
		if !valid {
			solved = db.Attempted
		}
		evaluation += comment + "<br>"
	}

	return evaluation, solved
}

// unzipBytes extracts a zip-from-memory into destDir.
// It will create destDir if it does not exist.
func unzipBytes(zipData []byte, destDir string) error {
	readerAt := bytes.NewReader(zipData)
	zr, err := zip.NewReader(readerAt, int64(len(zipData)))
	if err != nil {
		return err
	}

	for _, f := range zr.File {
		targetPath := filepath.Join(destDir, f.Name)

		// SECURITY: block Zip Slip (`../` escape)
		if !strings.HasPrefix(filepath.Clean(targetPath), filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path in zip: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		outFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}

		if _, err := io.Copy(outFile, rc); err != nil {
			outFile.Close()
			return err
		}

		if err := outFile.Close(); err != nil {
			return err
		}
	}

	return nil
}
