package main

import "fmt"

func main() {
	fmt.Println(begrüßen("Guten Morgen", "Judge"))
	// Ergebnis: GUTEN MORGEN, Judge

	fmt.Println(rahmen("Judge!", 8))
	// Ergebnis:
	// ********
	// Judge!
	// ********

	fmt.Println(aufräumen(`
		************************************
		*    JETZT KAUFEN UND 10% SPAREN   *
		************************************
	`))
	// Ergebnis: JETZT KAUFEN UND 10% SPAREN
}

// TODO: Definiere eine Funktion 'begrüßen'.
// - Die Funktion nimmt zwei Parameter:
//		- Die Grußformel
//		- Den Namen
// - Die Funktion soll die Grußformel in Großbuchstaben
//	 gefolgt von einem Komma und dann gefolgt vom Namen zurückgeben

// TODO: Definiere eine Funktion 'rahmen'
// - Die Funktion nimmt zwei Parameter:
//		- Den Namen
//		- Die Anzahl an Sternen für den Rahmen
// - Die Funktion soll den Namen mit einem
//	 Rahmen von Sternen darüber und darunter Zeichnen
//

// TODO: Definiere eine Funktion 'aufräumen'
// - Die Funktion akzeptiert einen Parameter, wie viele Autos gebaut werden sollen
// - Jedes Auto kostet 10.000€ aber 10 Autos können für 95.000€ produziert werden
//   Beispiel:
//		37 Autos können auf folgende Weise produziert werden:
//			37 = 3 Gruppen von 10 + 7 individuelle Autos
//			   = 3 * 95.000 + 7 * 10.000
//			   = 355.000
// - Tipp: Versucht die Funktion zu finden zum Ersetzen von Zeichen in:
//	 https://pkg.go.dev/strings
// - Tipp: Du brauchst neben obiger Funktion ebenfalls eine Funktion um Leerzeichen zu entfernen
