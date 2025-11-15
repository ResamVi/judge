package cli

import (
	"archive/zip"
	"bytes"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"io"
	"net/http"
	"os"
	"time"
)

type downloadResult struct {
	Error string
}

func updateDownload(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.ListExercises, cmd = m.ListExercises.Update(msg)
	cmds = append(cmds, cmd)

	m.Spinner, cmd = m.Spinner.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ListExercises.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			if m.Validating {
				return m, tea.Batch(cmds...)
			}

			i, ok := m.ListExercises.SelectedItem().(item)
			if !ok {
				panic("not found")
			}
			m.Folder = string(i)
			m.Validating = true

			cmds = append(cmds,
				m.Spinner.Tick,
				downloadCmd(m.Token, string(i)),
			)
		}
	case downloadResult:
		m.Validating = false
		if msg.Error == "" {
			m.Page++
			m.FarewellMessage = "Übung erfolgreich heruntergeladen"
			m.Quitting = true
			return m, nil
		}
		m.ErrorMessage = msg.Error
		m.Folder = ""
	}

	return m, tea.Batch(cmds...)
}

func viewDownload(m Model) string {
	tpl := m.ListExercises.View() + "\n"

	if m.Validating {
		tpl += fmt.Sprintf(" %s Lade herunter...\n\n", m.Spinner.View())
	} else if m.ErrorMessage != "" {
		tpl += fmt.Sprintf(" 💩 Fehlgeschlagen (%s)\n\n", m.ErrorMessage)
	}

	tpl += subtleStyle.Render("up/down: auswählen") + dotStyle +
		subtleStyle.Render("enter: bestätigen") + dotStyle +
		subtleStyle.Render("strg+c: beenden")

	return tpl
}

func downloadCmd(token, exercise string) tea.Cmd {
	return func() tea.Msg {
		return downloadResult{Error: downloadExercise(token, exercise)}
	}
}

func downloadExercise(token, exercise string) string {
	time.Sleep(1 * time.Second)

	fileName := Exercises[exercise]

	err := os.Mkdir(fileName, os.ModePerm)
	if err != nil {
		return "Ordner existiert schon"
	}

	req, err := http.NewRequest(http.MethodGet, JudgeURL+"/exercises/"+fileName+"/code", nil)
	if err != nil {
		return "NewRequest: " + err.Error()
	}
	req.Header.Set("token", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "Do: " + err.Error()
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "Unbekannter Status: " + resp.Status
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "ReadAll: " + err.Error()
	}

	// Create the file
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return "zip.NewReader: " + err.Error()
	}

	// Read all the files from zip archive
	for _, zipFile := range zipReader.File {
		unzippedFileBytes, err := readZipFile(zipFile)
		if err != nil {
			return "readZipFile: " + zipFile.Name
		}

		out, err := os.Create(fileName + "/" + zipFile.Name)
		if err != nil {
			return "os.Create: " + err.Error()
		}
		_, err = out.Write(unzippedFileBytes)
		if err != nil {
			return "out: " + err.Error()
		}

		out.Close()
	}

	return ""
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}
