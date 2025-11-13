package cli

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

func viewChoices(m Model) string {
	c := m.Choice

	tpl := "Hi! Was möchtest du tun?\n\n"
	tpl += "%s\n\n"
	tpl += subtleStyle.Render("up/down: auswählen") + dotStyle +
		subtleStyle.Render("enter: bestätigen") + dotStyle +
		subtleStyle.Render("esc: beenden")

	choices := fmt.Sprintf(
		"%s\n%s\n%s",
		checkbox("Neue Aufgabe herunterladen", c == 0),
		checkbox("Bearbeitete Aufgabe hochladen", c == 1),
		checkbox("Zugang konfigurieren", c == 2),
	)

	return fmt.Sprintf(tpl, choices)
}

func updateChoices(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.Choice++
			if m.Choice > 2 {
				m.Choice = 2
			}
		case "k", "up":
			m.Choice--
			if m.Choice < 0 {
				m.Choice = 0
			}
		case "enter":
			m.Chosen = true
			return m, frame()
		}
	}

	return m, nil
}
