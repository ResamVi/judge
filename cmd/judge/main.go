package main

import (
	"fmt"
	"github.com/ResamVi/judge/cli"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(cli.NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}
