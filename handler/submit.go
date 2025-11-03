package handler

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/fil
	"strings"
i
)

func (k Handler) Submit(c echo.Context) error {
	//token := c.Request().Header.Get("token")
	//exercise := c.Request().Header.Get("exercise")

	//user, err := k.db.GetUserFromToken(c.Request().Context(), token)
	//if err != nil {
	//	slog.Error("could not find user with token", "token", token, "err", err)
	//	return c.String(http.StatusNotFound, err.Error())
	//}

	// Decode base64
	data, err := io.ReadAll(base64.NewDecoder(base64.StdEncoding, c.Request().Body))
	if err != nil {
		slog.Error("bad base64", "error", err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	// Destination path
	destDir := filepath.Join("submissions", time.Now().Format("2006-01-02T15-04"))

	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		slog.Error("cannot create dest dir", "error", err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := unzipBytes(data, destDir); err != nil {
		slog.Error("unzip failed", "error", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, ":)")
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
