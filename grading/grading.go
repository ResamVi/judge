package grading

import (
	"regexp"
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

var Grading = map[string]Exercise{
	"01-compiler": {
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
}

func OutputMatches(expected string) Criteria {
	return Criteria{
		Description: "Ausgabe des Programms ist wie erwartet",
		Valid: func(code, output string) (string, bool) {
			if expected == output {
				return "✅ Ausgabe des Programms ist wie erwartet", true
			}

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
