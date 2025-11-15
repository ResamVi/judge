package main

import (
	"flag"
	"fmt"
	"github.com/ResamVi/judge/cli"
	tea "github.com/charmbracelet/bubbletea"
)

var dev = flag.Bool("dev", false, "development mode")

func main() {
	flag.Parse()

	p := tea.NewProgram(cli.NewModel(*dev))
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}
