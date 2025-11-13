package cli

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	// spring green-ish
	SuccessStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00ff7f")).MarginLeft(2)
)

type uploadResult struct {
	Error string
}

func updateUpload(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.List, cmd = m.List.Update(msg)
	cmds = append(cmds, cmd)

	m.Spinner, cmd = m.Spinner.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			if m.Validating {
				return m, tea.Batch(cmds...)
			}

			i, ok := m.List.SelectedItem().(item)
			if !ok {
				panic("not found")
			}
			m.Folder = string(i)
			m.Validating = true

			cmds = append(cmds,
				m.Spinner.Tick,
				uploadCmd(m.JudgeURL, m.Token, string(i)),
			)
		}
	case uploadResult:
		m.Validating = false
		if msg.Error == "" {
			m.Page++
			m.Quitting = true
			return m, nil
		}
		m.ErrorMessage = msg.Error
		m.Folder = ""
	}

	return m, tea.Batch(cmds...)
}

func viewUpload(m Model) string {
	tpl := m.List.View() + "\n"

	if m.Validating {
		tpl += fmt.Sprintf(" %s Lade hoch...\n\n", m.Spinner.View())
	} else if m.ErrorMessage != "" {
		tpl += fmt.Sprintf(" 💩 Fehlgeschlagen (%s)\n\n", m.ErrorMessage)
	}

	tpl += subtleStyle.Render("up/down: auswählen") + dotStyle +
		subtleStyle.Render("enter: bestätigen") + dotStyle +
		subtleStyle.Render("strg+c: beenden")

	return tpl
}

func uploadCmd(url, token, value string) tea.Cmd {
	return func() tea.Msg {
		return uploadResult{Error: uploadFolder(url, token, value)}
	}
}

func uploadFolder(url, token, directoryName string) string {
	time.Sleep(1 * time.Second)

	path, err := os.Getwd()
	if err != nil {
		return "Getwd: " + err.Error()
	}

	directory := filepath.Join(path, directoryName)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return "os.Stat: " + err.Error()
	}

	// 1. Zip the folder
	zipBytes, err := zipDirectory(directory)
	if err != nil {
		return "zipDirectory: " + err.Error()
	}

	// 2. Base64-encode zip
	b64 := base64.StdEncoding.EncodeToString(zipBytes)

	// 4. POST it
	req, err := http.NewRequest(http.MethodPost, url+"/submission", bytes.NewReader([]byte(b64)))
	if err != nil {
		return "NewRequest: " + err.Error()
	}
	req.Header.Set("token", token)
	req.Header.Set("exercise", directoryName)
	req.Header.Set("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "Do: " + err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "token unbekannt"
	}

	if resp.StatusCode == http.StatusOK {
		return ""
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "ReadAll: " + err.Error()
	}

	return string(b)
}

// zipDirectory walks dirPath and writes a .zip (in-memory) with full structure.
// Returns the raw bytes of the zip file.
func zipDirectory(dirPath string) ([]byte, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		// Build the path inside the zip (relative to dirPath)
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		// Skip the root itself ("."), but still include its children
		if relPath == "." {
			return nil
		}

		// Directories in zip need a trailing slash
		if info.IsDir() {
			_, err := zipWriter.Create(relPath + "/")
			return err
		}

		// It's a file: create header, copy contents
		fileHeader, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		fileHeader.Name = relPath
		fileHeader.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(fileHeader)
		if err != nil {
			return err
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(writer, f)
		return err
	})

	if err != nil {
		zipWriter.Close()
		return nil, err
	}

	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
