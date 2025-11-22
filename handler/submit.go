package handler

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/ResamVi/judge/db"
	"github.com/ResamVi/judge/grading"
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

	code, err := getCode(data)
	if err != nil {
		slog.Error("failed to get code", "error", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	destDir := filepath.Join("submissions", time.Now().Format("2006-01-02T15-04")+"_"+user.Username)

	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		slog.Error("cannot create dest dir", "error", err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := unzipBytes(data, destDir); err != nil {
		slog.Error("unzip failed", "error", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// == Build & run code ==
	go k.executeSubmission(user, destDir, exercise, code)

	slog.Info("submission received", "user", user.Username, "exercise", exercise)
	return c.NoContent(http.StatusOK)
}

func getCode(data []byte) (string, error) {
	readerAt := bytes.NewReader(data)
	zr, err := zip.NewReader(readerAt, int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("bad zip data: %w", err)
	}

	code := ""
	for _, f := range zr.File {
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

	return code, nil
}

func (k Handler) executeSubmission(user db.User, destDir, exercise, code string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//err := exec.Command("bash", "-c", fmt.Sprintf("cd %s && go mod tidy", destDir)).Run()
	//if err != nil {
	//	slog.Error("go mod tidy failed", "username", user.Username, "userId", user.ID, "exercise", exercise, "error", err.Error())
	//}

	runCmd := exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf("cd %s && go mod tidy && go run main.go", destDir))

	var buildStderr, buildStdout bytes.Buffer
	runCmd.Stderr = &buildStderr
	runCmd.Stdout = &buildStdout

	if f, ok := grading.Lazy[exercise]; ok {
		go f(runCmd)
		time.Sleep(1 * time.Second) // TODO: lol
	}

	var evaluation string
	var solved grading.Grade
	var output string

	if err := runCmd.Run(); err != nil {
		slog.Warn("submission returned errors", "username", user.Username, "output", buildStderr.String(), "error", err)
		solved = grading.Attempted

		if errors.Is(err, context.Canceled) {
			output = buildStderr.String()
			evaluation = "❌ Programm hat länger als 10 Sekunden gebraucht"
		} else {
			output = buildStderr.String()
			evaluation = "❌ Programm konnte nicht kompiliert werden"
		}
	} else {
		output = buildStdout.String()
		evaluation, solved = grading.GradeSubmission(exercise, code, output)
	}

	err := k.db.CreateSubmission(context.Background(), db.CreateSubmissionParams{
		UserID:     user.ID,
		ExerciseID: exercise,
		Code:       code,
		Output:     output,
		Evaluation: evaluation,
		Solved:     int32(solved),
	})
	if err != nil {
		slog.Error("failed to create submission", "username", user.Username, "userId", user.ID, "exercise", exercise, "error", err.Error())
	}

	err = os.RemoveAll(destDir)
	if err != nil {
		slog.Error("failed to remove folder", "destDir", destDir, "error", err.Error())
	}
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
