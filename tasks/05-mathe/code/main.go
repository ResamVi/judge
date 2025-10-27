package main

import "fmt"

func main() {
	fmt.Println(berechneRate(1547, 90))
	fmt.Println(berechneRateMinütlich(4))
	fmt.Println(berechneKosten(21))
}

// TODO: Definiere eine Funktion 'berechneRate'.
// - Die Funktion nimmt die Anzahl von Autos die produziert
//	 werden pro Stunde und die Erfolgschance und berechnet wie viele Autos
//	 erfolgreich in der Stunde produziert werden
// - Beide Parameter sind vom Typ 'int'. Die Erfolgschance wird als Zahl zwischen 0 und 100 gegeben
// - Der Rückgabewert soll vom Typ 'float64' sein

// TODO: Definiere eine Funktion 'berechneRateMinütlich'
// - Die Funktion macht das gleiche wie 'berechneRate' aber gibt zurück wie viele Autos *minütlich* produziert werden
// - Beide Parameter sind vom Typ 'int'. Die Erfolgschance wird als Zahl zwischen 0 und 100 gegeben
// - Der Rückgabewert soll vom Typ 'int' sein

// TODO: Definiere eine Funktion 'berechneKosten'
// - Die Funktion akzeptiert einen Parameter, wie viele Autos gebaut werden sollen
// - Jedes Auto kostet 10.000€ aber 10 Autos können für 95.000€ produziert werden
// Beispiel:
//		37 Autos können auf folgende Weise produziert werden:
//			37 = 3 Gruppen von 10 + 7 individuelle Autos
//			   = 3 * 95.000 + 7 * 10.000
//			   = 355.000
