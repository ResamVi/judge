package cli

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kirsle/configdir"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
)

const (
	dotChar  = " • "
	JudgeURL = "http://localhost:8080"
)

var (
	// Content should not be at the left most border
	mainStyle = lipgloss.NewStyle().MarginLeft(2)

	// Help text
	subtleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	// Checkbox style
	checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	// Dark grey separators for the help text using a dot
	dotStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)

	// A blue spinner
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)

type Model struct {
	// Shared
	Page            int
	Validating      bool
	ErrorMessage    string
	Quitting        bool
	JudgeURL        string
	FarewellMessage string

	// First View
	Token      string
	TokenInput textinput.Model
	Spinner    spinner.Model

	// Second View
	Choice int // Which command (herunterladen, hochladen, konfigurieren)

	// Third View (upload)
	Folder      string
	ListFolders list.Model

	// Third View (download)
	ListExercises list.Model
}

var (
	titleStyle = lipgloss.NewStyle()
	Exercises  map[string]string
)

func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "7e41c20a-37df-4dd5-b43c-d112a0f9dc1f"
	ti.Prompt = "Token: "
	ti.CharLimit = 36
	ti.Width = 36
	ti.Focus()

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	// List of folders in current working directory
	entries, err := os.ReadDir(".")
	if err != nil {
		panic(err)
	}

	var folders []list.Item
	for _, e := range entries {
		if e.Type() == fs.ModeDir {
			folders = append(folders, item(e.Name()))
		}
	}

	l := list.New(folders, itemDelegate{}, 38, len(folders)*3)
	l.Title = "Wähle den Ordner zum Hochladen aus"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.Styles.Title = titleStyle

	// List of exercises of judge
	resp, err := http.Get(JudgeURL + "/exercises")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&Exercises)
	if err != nil {
		panic(err)
	}
	var rawExercises []string
	for k, _ := range Exercises {
		rawExercises = append(rawExercises, k)
	}
	sort.Strings(rawExercises)

	var exercises []list.Item
	for _, ex := range rawExercises {
		exercises = append(exercises, item(ex))
	}

	ll := list.New(exercises, itemDelegate{}, 45, 8)
	ll.Title = "Wähle die Aufgabe aus"
	ll.SetShowStatusBar(false)
	ll.SetFilteringEnabled(false)
	ll.SetShowHelp(false)
	ll.Styles.Title = titleStyle

	m := Model{
		TokenInput:    ti,
		Spinner:       s,
		ListFolders:   l,
		ListExercises: ll,
	}

	// Usage should be as easy as possible.
	// Therefore, check whether the user has already used the judge
	// and re-use the token. Otherwise, ask for the user to provide it first.
	configFile := filepath.Join(configdir.LocalConfig("judge"), "settings.json")

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return m
	}

	fh, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	type Config struct {
		Token string `json:"token"`
	}
	var config Config
	err = json.NewDecoder(fh).Decode(&config)
	if err != nil {
		panic(err)
	}
	m.Token = config.Token
	m.Page = 1

	return m
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// When user wants to quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}

	switch m.Page {
	case 0:
		return updateAuthenticate(msg, m)
	case 1:
		return updateChoices(msg, m)
	case 2:
		switch m.Choice {
		case 0:
			return updateUpload(msg, m)
		case 1:
			return updateDownload(msg, m)
		case 2:
			m.Page = 0
			m.Choice = 0
			return m, nil
		}
	}

	return m, tea.Quit
}

func (m Model) View() string {
	if m.Quitting {
		return SuccessStyle.Render("\n" + m.FarewellMessage + "\n\n")
	}

	var s string
	switch m.Page {
	case 0:
		s = viewAuthenticate(m)
	case 1:
		s = viewChoices(m)
	case 2:
		switch m.Choice {
		case 0:
			s = viewUpload(m)
		case 1:
			s = viewDownload(m)
		case 2:
			s = viewAuthenticate(m)
		}
	}

	return mainStyle.Render("\n" + s + "\n\n") // Center content vertically
}

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}

func hyperlink(url, text string) string {
	return fmt.Sprintf(
		"\x1B]8;;%s\x1B\\%s\x1B]8;;\x1B\\",
		url,
		text,
	)
}
