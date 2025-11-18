package grading

import (
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"regexp"
	"sort"
	"strings"
)

type Grade int

const (
	NotAttempted Grade = 0
	Attempted    Grade = 1
	Solved       Grade = 2
)

type Exercise struct {
	Criteria []Criteria
}

type Criteria struct {
	Description string
	Valid       func(code, output string) (comment string, failed bool)
}

var Lazy = map[string]func(cmd *exec.Cmd){
	"01-judge-einrichten":        func(cmd *exec.Cmd) {},
	"02-hello-world":             func(cmd *exec.Cmd) {},
	"03-werte-ausgeben":          func(cmd *exec.Cmd) {},
	"04a-variablen-kennenlernen": func(cmd *exec.Cmd) {},
	"04b-variablen-tauschen":     func(cmd *exec.Cmd) {},
	"05-arithmetik": func(cmd *exec.Cmd) {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			slog.Error("failed to get stdin pipe", "error", err.Error())
		}
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, "4\n")
			io.WriteString(stdin, "4\n")
			io.WriteString(stdin, "5\n")
		}()

	},
	"X-hacking": func(cmd *exec.Cmd) {},
}

var Grading = map[string]Exercise{
	"01-judge-einrichten": {
		Criteria: []Criteria{
			{
				Description: "Programm blieb unverändert",
				Valid: func(code, output string) (string, bool) {
					if code == "package main\n\nimport (\n\t\"fmt\"\n)\n\nfunc main() {\n\tlines := [][]rune{\n\t\t{32, 47, 92, 95, 47, 92},\n\t\t{40, 32, 111, 46, 111, 32, 41},\n\t\t{32, 62, 32, 94, 32, 60},\n\t}\n\n\tfor _, line := range lines {\n\t\tfmt.Println(string(line))\n\t}\n}\n" {
						return "✅ Programmcode blieb unverändert", true
					}
					return "❌ Programm wurde verändert", false
				},
			},
		},
	},
	"02-hello-world": {
		Criteria: []Criteria{
			NoHackingAttempt,
			OutputMatches("Hello World!\n"),
		},
	},
	"03-werte-ausgeben": {
		Criteria: []Criteria{
			NoHackingAttempt,
			OutputMatches("42 3.141 Go macht Spaß true\n"),
			CodeWithout(`"42"`, "Zahlen brauchen keine Anführungsstriche"),
			CodeWithout(`"3.141"`, "Kommazahlen brauchen keine Anführungsstriche"),
			CodeWithout(`"true"`, "Wahrheitswerte brauchen keine Anführungsstriche"),
			CodeWithout(`"42 3.141 Go macht Spaß true"`, "Benutze die Möglichkeit mehrere Werte in fmt.Println zu packen"),
		},
	},
	"04a-variablen-kennenlernen": {
		Criteria: []Criteria{
			NoHackingAttempt,
			CodeRegex(`\w+ :=`, "Kurze Variablen Deklaration vorhanden"),
			CodeRegex(`var \w+`, "Normale Variablendeklaration vorhanden"),
			CodeRegex(`\w+, \w+`, "Mehrere Variablen gleichzeitig wurden deklariert"),
		},
	},
	"04b-variablen-tauschen": {
		Criteria: []Criteria{
			NoHackingAttempt,
			OutputMatches("27 5\n"),
			CodeRegexNot(`=.+\d+`, "Variablen dürfen nicht einfach anderen Zahlen zugewiesen werden", 5),
			CodeRegex(`fmt\.Println\(a\, b\)`, "fmt.Println wurde nicht verändert"),
		},
	},

	"05-arithmetik": {
		Criteria: []Criteria{
			NoHackingAttempt,
			OutputMatches(`Was ist deine Note in Geometrie?
Was ist deine Note in Algebra?
Was ist deine Note in Physik?
Dein Notendurchschnitt:
4.3333335
false
`),
		},
	},

	"X-hacking": {
		Criteria: []Criteria{
			NoHackingAttempt,
			{
				Description: "Solange man das liest funktioniert alles noch",
				Valid: func(code, output string) (string, bool) {
					return "❌ Der Server steht noch", false
				},
			},
		},
	},
}

func CodeWithout(avoid string, explanation string) Criteria {
	return Criteria{
		Description: "Ausgabe des Programms ist wie erwartet",
		Valid: func(code, output string) (string, bool) {
			if !strings.Contains(code, avoid) {
				return "✅ Programm vermeidet Fehler (<i>" + explanation + "</i>)", true
			}

			return "❌ Program enthält unerwünschten Code (<i>" + explanation + "</i>)", false
		},
	}
}

func CodeRegexNot(expected string, explanation string, ignored ...int) Criteria {
	return Criteria{
		Description: "Ausgabe des Programms ist wie erwartet",
		Valid: func(code, output string) (string, bool) {
			code = removeLines(code, ignored)

			if !regexp.MustCompile(expected).MatchString(code) {
				return "✅ Programm erfüllt Anforderung (<i>" + explanation + "</i>)", true
			}

			return "❌ Program erfüllt Anforderung nicht (<i>" + explanation + "</i>)", false
		},
	}
}

func CodeRegex(expected string, explanation string, ignored ...int) Criteria {
	return Criteria{
		Description: "Ausgabe des Programms ist wie erwartet",
		Valid: func(code, output string) (string, bool) {
			code = removeLines(code, ignored)

			if regexp.MustCompile(expected).MatchString(code) {
				return "✅ Programm erfüllt Anforderung (<i>" + explanation + "</i>)", true
			}

			return "❌ Program erfüllt Anforderung nicht (<i>" + explanation + "</i>)", false
		},
	}
}

func OutputMatches(expected string) Criteria {
	return Criteria{
		Description: "Ausgabe des Programms ist wie erwartet",
		Valid: func(code, output string) (string, bool) {
			if expected == output {
				return "✅ Ausgabe des Programms ist wie erwartet", true
			}
			slog.Warn("mismatch: ", "output", output, "expected", expected)

			return "❌ Ausgabe des Programms ist nicht wie erwartet:<br><pre><code>" + expected + "</code></pre>", false
		},
	}
}

var NoHackingAttempt = Criteria{
	Description: "Hat keine unzulässigen Systemzugriffe",
	Valid: func(code, output string) (string, bool) {
		patterns := []string{
			`(?i)(subprocess|exec\.|shell|eval|child_process)`, // any shell execution commands to spawn new processes
			`(?i)("os")`,                         // operating system operations can stop program/create big files/read filesystem
			`(?i)(net\.Listen|net\.Dial|http\.)`, // net/http calls can communicate with remote servers
		}
		for _, pattern := range patterns {
			if regexp.MustCompile(pattern).MatchString(code) {
				return "❌ Unzulässiger Systemzugriff erkannt", false
			}
		}
		return "✅ Code hat keine unzulässigen Systemzugriffe", true
	},
}

func GradeSubmission(exercise string, code string, output string) (string, Grade) {
	fmt.Println(output)

	criteria, ok := Grading[exercise]
	if !ok {
		return "Unbekannt: " + exercise, NotAttempted
	}

	evaluation := "✅ Programm konnte kompiliert werden<br>" // Already checked when cmd was executed
	solved := Solved

	for _, fn := range criteria.Criteria {
		comment, valid := fn.Valid(code, output)
		if !valid {
			solved = Attempted
		}
		evaluation += comment + "<br>"
	}

	return evaluation, solved
}

func removeLines(code string, indices []int) string {
	lines := strings.Split(code, "\n")

	// Sort indices descending so removal does not shift early indexes
	sort.Sort(sort.Reverse(sort.IntSlice(indices)))

	for _, n := range indices {
		if n < 0 || n >= len(lines) {
			continue // ignore invalid
		}
		lines = append(lines[:n], lines[n+1:]...)
	}

	return strings.Join(lines, "\n")
}
