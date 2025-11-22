package grading

import (
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

// Lazy because I cba to do this cleaner
var Lazy = map[string]func(cmd *exec.Cmd){
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

	"10a-for-schleife": func(cmd *exec.Cmd) {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			slog.Error("failed to get stdin pipe", "error", err.Error())
		}
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, "13\n")
		}()
	},
	"10b-fakultaet": func(cmd *exec.Cmd) {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			slog.Error("failed to get stdin pipe", "error", err.Error())
		}
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, "14\n")
		}()
	},
	"10c-durchschnitt-berechnen": func(cmd *exec.Cmd) {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			slog.Error("failed to get stdin pipe", "error", err.Error())
		}
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, "3\n")
			io.WriteString(stdin, "7\n")
			io.WriteString(stdin, "5\n")
			io.WriteString(stdin, "-5\n")
			io.WriteString(stdin, "0\n")
		}()
	},
	"10d-zahl-erraten": func(cmd *exec.Cmd) {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			slog.Error("failed to get stdin pipe", "error", err.Error())
		}
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, "0\n")
			io.WriteString(stdin, "1\n")
			io.WriteString(stdin, "2\n")
			io.WriteString(stdin, "3\n")
			io.WriteString(stdin, "4\n")
			io.WriteString(stdin, "5\n")
			io.WriteString(stdin, "6\n")
			io.WriteString(stdin, "7\n")
			io.WriteString(stdin, "8\n")
			io.WriteString(stdin, "9\n")
		}()
	},
	"14-errors": func(cmd *exec.Cmd) {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			slog.Error("failed to get stdin pipe", "error", err.Error())
		}
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, "023\n")
		}()
	},
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

	"06-funktionen": {
		Criteria: []Criteria{
			NoHackingAttempt,
			OutputMatches("690\n830\n460\n"),
			CodeRegex(`fmt.Println\(berechneGehalt\(10\, 3\)\)`, "fmt.Println(10, 3) wurde nicht verändert"),
			CodeRegex(`fmt.Println\(berechneGehalt\(20\, 1\)\)`, "fmt.Println(20, 1) wurde nicht verändert"),
			CodeRegex(`fmt.Println\(berechneGehalt\(3\, 0\)\)`, "fmt.Println(3, 0) wurde nicht verändert"),
		},
	},

	"07-booleans": {
		Criteria: []Criteria{
			NoHackingAttempt,
			OutputMatches(`Hat ein Ticket: true
Ist VIP: true
Hat weder Ticket noch VIP: false
Nicht eingesteckt: false
Nicht angeschalten: false
Eingesteckt und angeschalten: true
Name enthält Zahlen: false
Name enthält keine Zahlen: true
101 Grad: true
100 Grad: false
99 Grad: false
`),
			CodeWith(`fmt.Printf("Hat ein Ticket: %v\n", eintrittErlaubt(true, false))`, "fmt.Printf wurde nicht verändert"),
			CodeWith(`fmt.Printf("Ist VIP: %v\n", eintrittErlaubt(false, true))`, "fmt.Printf wurde nicht verändert"),
			CodeWith(`fmt.Printf("Hat weder Ticket noch VIP: %v\n", eintrittErlaubt(false, false))`, "fmt.Printf wurde nicht verändert"),
			CodeWith(`fmt.Printf("Nicht eingesteckt: %v\n", computerLäuft(false, true))`, "fmt.Printf wurde nicht verändert"),
			CodeWith(`fmt.Printf("Nicht angeschalten: %v\n", computerLäuft(true, false))`, "fmt.Printf wurde nicht verändert"),
			CodeWith(`fmt.Printf("Eingesteckt und angeschalten: %v\n", computerLäuft(true, true))`, "fmt.Printf wurde nicht verändert"),
			CodeWith(`fmt.Printf("Name enthält Zahlen: %v\n", nameValide(true))`, "fmt.Printf wurde nicht verändert"),
			CodeWith(`fmt.Printf("Name enthält keine Zahlen: %v\n", nameValide(false))`, "fmt.Printf wurde nicht verändert"),
			CodeWith(`fmt.Printf("101 Grad: %v\n", istHeiß(101))`, "fmt.Print wurde nicht verändert"),
			CodeWith(`fmt.Printf("100 Grad: %v\n", istHeiß(100))`, "fmt.Printf wurde nicht verändert"),
			CodeWith(`fmt.Printf("99 Grad: %v\n", istHeiß(99))`, "fmt.Printf wurde nicht verändert"),
		},
	},

	"08-if-bedingung": {
		Criteria: []Criteria{
			NoHackingAttempt,
			OutputMatches("true\nfalse\ntrue\n800\n700\n500\n"),
			CodeWith(`func main() {
	fmt.Println(brauchtFührerschein("auto"))
	fmt.Println(brauchtFührerschein("fahrrad"))
	fmt.Println(brauchtFührerschein("lkw"))

	fmt.Println(schätzeWert(1000, 1))
	fmt.Println(schätzeWert(1000, 5))
	fmt.Println(schätzeWert(1000, 15))
}`, "Die main Funktion wurde nicht verändert"),
		},
	},

	"09-switch": {
		Criteria: []Criteria{
			NoHackingAttempt,
			OutputMatches("11\n10\n3\n10\n10\n10\n0\nTeilen\nSieg\nKarte\nStehen\nKarte\nStehen\n"),
			CodeWith(`func main() {
	fmt.Println(kartenwert("Ass"))
	fmt.Println(kartenwert("Zehn"))
	fmt.Println(kartenwert("Drei"))
	fmt.Println(kartenwert("König"))
	fmt.Println(kartenwert("Bube"))
	fmt.Println(kartenwert("Dame"))
	fmt.Println(kartenwert("unbekannt"))

	fmt.Println(ersterZug("Ass", "Ass", "Zehn"))
	fmt.Println(ersterZug("Zehn", "Ass", "Neun"))
	fmt.Println(ersterZug("Zehn", "Ass", "Zehn"))
	fmt.Println(ersterZug("Zehn", "Sieben", "Zehn"))
	fmt.Println(ersterZug("Zehn", "Sechs", "Zehn"))
	fmt.Println(ersterZug("Zehn", "Sechs", "Drei"))
}`, "Die main Funktion wurde nicht verändert"),
		},
	},

	"10a-for-schleife": {
		Criteria: []Criteria{
			NoHackingAttempt,
			OutputMatches("Nenne eine Zahl:\n1\n2\n3\n4\n5\n6\n7\n8\n9\n10\n11\n12\n13\n"),
		},
	},

	"10b-fakultaet": {
		Criteria: []Criteria{
			NoHackingAttempt,
			OutputMatches("Nenne eine Zahl:\n87178291200\n"),
		},
	},

	"10c-durchschnitt-berechnen": {
		Criteria: []Criteria{
			NoHackingAttempt,
			OutputMatches("Gebe mir Zahlen:\n5\n"),
		},
	},

	"10d-zahl-erraten": {
		Criteria: []Criteria{
			NoHackingAttempt,
			CodeWith("rand.Intn(", "Sollte das math/rand Package benutzen"),
			OutputRegex("Versuche die Zahl zu erraten:", "Frage am Anfang: 'Versuche die Zahl zu erraten:'"),
			OutputRegex(`Anzahl der Versuche: \d+`, "Antwort am Ende: 'Anzahl der Versuche: N'"),
		},
	},

	"10e-fizz-buzz": {
		Criteria: []Criteria{
			NoHackingAttempt,
			OutputMatches("1\n2\nFizz\n4\nBuzz\nFizz\n7\n8\nFizz\nBuzz\n11\nFizz\n13\n14\nFizzBuzz\n16\n17\nFizz\n19\nBuzz\n"),
		},
	},

	"11-slices": {
		Criteria: []Criteria{
			NoHackingAttempt,
			OutputMatches("[1 2 3 4 5 6 7 8 9 10 J Q K A]\n4\n[1 2 3 4 5 6 7 8 9 10 J Q K K]\n[1 2 3 4 5 6 7 8 9 10 J Q K A 1]\n"),
			CodeWith(`func main() {
	fmt.Println(stapel())
	fmt.Println(karteNehmen(stapel(), 3))
	fmt.Println(karteAustauschen(stapel(), 13, "K"))
	fmt.Println(karteHinzufügen(stapel(), "1"))
}`, "Die main Funktion wurde nicht verändert"),
		},
	},

	"12a-structs": {
		Criteria: []Criteria{
			NoHackingAttempt,
			OutputRegex("60", "Fläche des Rechtecks wird berechnet"),
			OutputRegex("Der Film \\w+ kam im Jahr \\d+ und gehört zur \\w+ Genre", "Werbung angezeigt"),
			OutputRegex(`main.Konto{Name:"Svenja Schmidt", Guthaben:100}`, "Konto angezeigt"),
			CodeWith(`fmt.Printf("%#v\n", kontoÖffnen("Svenja Schmidt", 100))`, "kontoÖffnen blieb unverändert"),
			CodeWith(`werbungZeigen(Film{`, "werbungZeigen kriegt ein Film struct übergeben"),
			CodeWith(`fmt.Println(flächeBerechnen(Rechteck{Höhe: 12, Breite: 5}))`, "flächeBerechnen blieb unverändert"),
		},
	},

	"12b-autorennen": {
		Criteria: []Criteria{
			NoHackingAttempt,
			CodeWith(`func main() {
	fmt.Println(baueAuto(50, 10))
	fmt.Println(fahren(baueAuto(50, 10)))
	fmt.Println(prüfeFahrt(Auto{batteriestand: 0}, Rennstrecke{distanz: 10}))
	fmt.Println(prüfeFahrt(Auto{batteriestand: 100, batterieverbrauch: 1, geschwindigkeit: 0}, Rennstrecke{distanz: 10}))
	fmt.Println(prüfeFahrt(Auto{batteriestand: 100, batterieverbrauch: 1, geschwindigkeit: 10}, Rennstrecke{distanz: 10}))
}`, "Die main Funktion wurde nicht verändert"),

			OutputMatches("{100 10 50 0}\n{90 10 50 50}\nfalse\nfalse\ntrue\n"),
		},
	},

	"13-module": {
		Criteria: []Criteria{
			NoHackingAttempt,
			CodeWith(`"github.com/sanity-io/litter"`, "Das Modul wurde importiert"),
			OutputMatches("main.Person{\n  Name: \"Max Mustermann\",\n  Alter: 40,\n  Adresse: main.Adresse{\n    Straße: \"Mustermannstraße\",\n    Stadt: \"Karlsruhe\",\n    PLZ: \"76131\",\n  },\n}\n"),
		},
	},

	"14-errors": {
		Criteria: []Criteria{
			NoHackingAttempt,
			CodeWith(`fmt.Println("Passwort angenommen")`, "Passwort angenommen wird ausgegeben"),
			OutputMatches("Gebe ein Passwort ein:\nPasswort ist zu kurz\n"),
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

func CodeWith(expected string, explanation string, ignored ...int) Criteria {
	return Criteria{
		Description: "Ausgabe des Programms ist wie erwartet",
		Valid: func(code, output string) (string, bool) {
			code = removeLines(code, ignored)

			if strings.Contains(code, expected) {
				return "✅ Programm erfüllt Anforderung (<i>" + explanation + "</i>)", true
			}

			return "❌ Program erfüllt Anforderung nicht (<i>" + explanation + "</i>)", false
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

func OutputRegex(expected string, explanation string) Criteria {
	return Criteria{
		Description: "Ausgabe des Programms ist wie erwartet",
		Valid: func(code, output string) (string, bool) {
			if regexp.MustCompile(expected).MatchString(output) {
				return "✅ Ausgabe des Programms ist wie erwartet (<i>" + explanation + "</i>)", true
			}
			slog.Warn("mismatch", "output", output, "expected", expected)

			return "❌ Ausgabe des Programms ist nicht wie erwartet (<i>" + explanation + "</i>)", false
		},
	}
}

var NoHackingAttempt = Criteria{
	Description: "Hat keine unzulässigen Systemzugriffe",
	Valid: func(code, output string) (string, bool) {
		patterns := []string{
			`(?i)subprocess`,                     // any shell execution commands to spawn new processes
			`(?i)exec\.`,                         // any shell execution commands to spawn new processes
			`(?i)shell`,                          // any shell execution commands to spawn new processes
			`eval`,                               // any shell execution commands to spawn new processes
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
