package main // 1. Jedes Programm fängt mit einem "package <name>" an

import "fmt" // 2. Um Text auszugeben brauchen wir das format (kurz: "fmt") Modul

// 3. Programmcode schreiben wir innerhalb Funktionen
// Die erste Funktion die von einem Programm ausgeführt wird ist "main"
func main() {
	// 4. Texte ausgeben mit einer Funktion aus dem fmt Modul
	fmt.Println("Hallo Welt")

	// 5. Typen kennengelernt
	var normaleZahl int
	var andereZahl int32
	var großeZahl int64

	var kommaZahl float32 = 3.0
	var großeKommazahl float64 = 4.0

	kurz := "Hallo" // Strings

	könnenSchweineFliegen := false // Booleans
	istDerPapstKatholisch := true

	fmt.Println(normaleZahl, andereZahl, großeZahl, kommaZahl, großeKommazahl, kurz, könnenSchweineFliegen, istDerPapstKatholisch)

	// 6. Wir können rechnen
	ergebnis := kommaZahl + kommaZahl
	fmt.Println(ergebnis) // Ausgabe: 6.0

	// 7. Wir können auch Zahlen vergleichen
	fmt.Println(ergebnis < 10) // true

	// 8. Wir können den Benutzer nach einer Eingabe fragen
	var name string
	fmt.Println("Wie ist dein Name?")
	fmt.Scanln(&name) // Das Programm wird so lange warten bis eine Eingabe und "Enter" gedrückt wurde

	// 9. Funktionen
	// neuesAlter := funktion(5)
	/*
		mehrzeilig
		mehrzeilig
		mehrzeilig
		mehrzeilig
	*/

}

func funktion(alter int) int {
	return alter + 1
}
