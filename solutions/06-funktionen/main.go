package main

import "fmt"

func main() {
	fmt.Println(berechneGehalt(10, 3))
	fmt.Println(berechneGehalt(20, 1))
	fmt.Println(berechneGehalt(3, 0))
}

// -- Ab hier editieren erlaubt --

// In einem Unternehmen berechnet sich das monatliche Gehalt eines Angestellten wie folgt:
// Mindestlohn 400€ im Monat,
// + zuzüglich 20€ multipliziert mit der Anzahl der Beschäftigungsjahre,
// + zuzüglich 30€ für jedes Kind.
//
// Aufgabe: Schreibe eine Funktion 'berechneGehalt' mit erstem Parameter die Beschäftigungsjahre
// und zweitem Parameter die Anzahl der Kinder und welches das Gehalt zurückgibt
func berechneGehalt(jahre int, kinder int) int {
	return 400 + 20*jahre + 30*kinder
}
