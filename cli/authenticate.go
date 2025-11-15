package cli

import (
	"encoding/json"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kirsle/configdir"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Token string `json:"token"`
}

type validationResult struct {
	Error string
}

func viewAuthenticate(m Model) string {
	tpl := "Hi! Du scheinst den Judge zum ersten Mal zu benutzen.\n"
	tpl += "Ich brauche deinen Token bevor du starten kannst.\n\n"
	tpl += "%s\n\n"

	if m.Validating {
		tpl += fmt.Sprintf(" %s Prüfe Token...\n\n", m.Spinner.View())
	} else if m.ErrorMessage != "" {
		tpl += fmt.Sprintf(" 💩 Token invalide (%s)\n\n", m.ErrorMessage)
	}

	tpl += subtleStyle.Render("Dein Token ist hier: ") + hyperlink("https://judge.resamvi.io/token", "https://judge.resamvi.io/token")
	tpl += subtleStyle.Render("\nstrg+c: beenden")

	return fmt.Sprintf(tpl, m.TokenInput.View())
}

func updateAuthenticate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.Spinner, cmd = m.Spinner.Update(msg)
	cmds = append(cmds, cmd)

	m.TokenInput, cmd = m.TokenInput.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.Validating {
				return m, tea.Batch(cmds...)
			}

			m.Validating = true
			m.Token = m.TokenInput.Value()
			m.TokenInput.Blur()

			cmds = append(cmds,
				m.Spinner.Tick,
				validateCmd(m.TokenInput.Value()),
			)
		}
	case validationResult:
		m.Validating = false
		m.ErrorMessage = msg.Error
		m.TokenInput.Focus()

		if msg.Error == "" {
			m.Page++
			return m, nil
		}
	}

	return m, tea.Batch(cmds...)
}
func validateCmd(value string) tea.Cmd {
	return func() tea.Msg {
		return validationResult{Error: validateToken(value)}
	}
}

func validateToken(token string) string {
	time.Sleep(1 * time.Second)

	req, err := http.NewRequest(http.MethodPost, JudgeURL+"/validate/token", nil)
	if err != nil {
		return "NewRequest: " + err.Error()
	}
	req.Header.Set("token", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "Do: " + err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "Token unbekannt"
	}

	if resp.StatusCode == http.StatusOK {
		configDir := configdir.LocalConfig("judge")
		configFile := filepath.Join(configDir, "settings.json")

		if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
			return "os.MkdirAll: " + err.Error()
		}

		fh, err := os.Create(configFile)
		if err != nil {
			return "os.Create: " + err.Error()
		}
		defer fh.Close()

		config := Config{Token: token}
		err = json.NewEncoder(fh).Encode(&config)
		if err != nil {
			return "Encode: " + err.Error()
		}

		return ""
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}
	return string(b)
}
