package cli

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fogleman/ease"
	"github.com/kirsle/configdir"
	"github.com/lucasb-eyer/go-colorful"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// TODO: Not needed maybe
const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
	dotChar           = " • "
)

var (
	// Content should not be at the left most border
	mainStyle = lipgloss.NewStyle().MarginLeft(2)

	keywordStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	progressEmpty = subtleStyle.Render(progressEmptyChar)

	// Dark grey separators for the help text
	dotStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)

	// Gradient colors we'll use for the progress bar
	ramp = makeRampStyles("#B14FFF", "#00FFA3", progressBarWidth)
)

type (
	tickMsg  struct{}
	frameMsg struct{}
)

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func frame() tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
		return frameMsg{}
	})
}

type Model struct {
	Page int // Which view

	// First View
	Validating   bool
	Token        string
	TextInput    textinput.Model
	Spinner      spinner.Model
	ErrorMessage string

	// Second View
	Choice int  // Which command (herunterladen, hochladen, konfigurieren)
	Chosen bool // User has chosen command

	Frames   int
	Progress float64
	Loaded   bool
	Quitting bool

	JudgeURL string
}

func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "7e41c20a-37df-4dd5-b43c-d112a0f9dc1f"
	ti.Prompt = "Token: "
	ti.CharLimit = 32
	ti.Width = 20
	ti.Focus()

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	m := Model{
		TextInput: ti,
		Spinner:   s,
		JudgeURL:  "http://localhost:8080",
	}

	// Usage should be as easy as possible.
	// Therefore, check whether the user has already used the judge
	// and re-use the token. Otherwise, ask for the user to provide it first.
	configFile := filepath.Join(configdir.LocalConfig("judge"), "settings.json")

	type Config struct {
		Token string `json:"token"`
	}
	var config Config

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return m
	}

	fh, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	err = json.NewDecoder(fh).Decode(&config)
	if err != nil {
		panic(err)
	}
	m.Token = config.Token
	m.Page = 1 // Skip authentication view (which is page 0)

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
	)
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

	if m.Page == 0 {
		return updateAuthenticate(msg, m)
	}

	// TODO: Switch
	if m.Page == 1 {
		return updateChoices(msg, m)
	}
	return updateChosen(msg, m)
}

func (m Model) View() string {
	if m.Quitting {
		return "" // Clear screen
	}

	var s string
	if m.Page == 0 {
		s = viewAuthenticate(m)
	}

	if m.Page == 1 {
		s = viewChoices(m)
	}
	//if {
	//	s = viewChosen(m)
	//}
	return mainStyle.Render("\n" + s + "\n\n") // Center content vertically
}

// Sub-update functions

// Update loop for the second view after a choice has been made
func updateChosen(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case frameMsg:
		if !m.Loaded {
			m.Frames++
			m.Progress = ease.OutBounce(float64(m.Frames) / float64(100))
			if m.Progress >= 1 {
				m.Progress = 1
				m.Loaded = true
				return m, tick()
			}
			return m, frame()
		}
	}

	return m, nil
}

// TODO: Check for token, endpoint that validates token

// The second view, after a task has been chosen
func viewChosen(m Model) string {
	var msg string

	switch m.Choice {
	case 0:
		msg = fmt.Sprintf("Carrot planting?\n\nCool, we'll need %s and %s...", keywordStyle.Render("libgarden"), keywordStyle.Render("vegeutils"))
	case 1:
		msg = fmt.Sprintf("A trip to the market?\n\nOkay, then we should install %s and %s...", keywordStyle.Render("marketkit"), keywordStyle.Render("libshopping"))
	case 2:
		msg = fmt.Sprintf("Reading time?\n\nOkay, cool, then we’ll need a library. Yes, an %s.", keywordStyle.Render("actual library"))
	default:
		msg = fmt.Sprintf("It’s always good to see friends.\n\nFetching %s and %s...", keywordStyle.Render("social-skills"), keywordStyle.Render("conversationutils"))
	}

	label := "Downloading..."
	if m.Loaded {
		label = fmt.Sprintf("Downloaded. Exiting in XXX seconds...")
	}

	return msg + "\n\n" + label + "\n" + progressbar(m.Progress) + "%"
}

func progressbar(percent float64) string {
	w := float64(progressBarWidth)

	fullSize := int(math.Round(w * percent))
	var fullCells string
	for i := 0; i < fullSize; i++ {
		fullCells += ramp[i].Render(progressFullChar)
	}

	emptySize := int(w) - fullSize
	emptyCells := strings.Repeat(progressEmpty, emptySize)

	return fmt.Sprintf("%s%s %3.0f", fullCells, emptyCells, math.Round(percent*100))
}

// Utils

// Generate a blend of colors.
func makeRampStyles(colorA, colorB string, steps float64) (s []lipgloss.Style) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0.0; i < steps; i++ {
		c := cA.BlendLuv(cB, i/steps)
		s = append(s, lipgloss.NewStyle().Foreground(lipgloss.Color(colorToHex(c))))
	}
	return
}

// Convert a colorful.Color to a hexadecimal format.
func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

// Helper function for converting colors to hex. Assumes a value between 0 and
// 1.
func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
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
