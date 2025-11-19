package main

import "fmt"

func main() {
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
}

// Definiere eine Funktion 'kartenwert'.
// - Die Funktion nimmt einen Parameter:
//   - Der Name der Karte
//
// - Die Funktion soll den Zahlenwert der Karte zurückgeben oder 0, falls unbekannt
//
// Folgende Zahlenwerte haben die Karten:
// Ass 		11 	Acht 	8
// Zwei 	2 	Neun 	9
// Drei 	3 	Zehn 	10
// Vier 	4 	Bube 	10
// Fünf 	5 	Dame 	10
// Sechs 	6 	König 	10
// Sieben 	7
func kartenwert(s string) int {
	switch s {
	case "Zwei":
		return 2
	case "Drei":
		return 3
	case "Vier":
		return 4
	case "Fünf":
		return 5
	case "Sechs":
		return 6
	case "Sieben":
		return 7
	case "Acht":
		return 8
	case "Neun":
		return 9
	case "Zehn", "Bube", "Dame", "König":
		return 10
	case "Ass":
		return 11

	}
	return 0
}

// Definiere eine Funktion 'ersterZug'
// - Die Funktion nimmt drei Parameter:
//   - Der Name deiner ersten Karte
//   - Der Name deiner zweiten Karte
//   - Der Name der Karte des Dealers, die offen ist
//
// - Die Funktion gibt den String zurück welcher Zug gemacht werden soll
//   - "Karte" (eine weitere Karte nehmen)
//   - "Stehen" (keine weitere mehr nehmen)
//   - "Teilen" (Karten aufteilen)
//   - "Sieg" (wenn man schon 21 hat)
//
// - Die Funktion soll folgende Strategie implementieren:
//   - Falls zwei Asse: "Teilen"
//   - Falls dein Kartenwert 21 ist und der Dealer kein Ass/10/Bube/Dame/König: Sieg (weil Blackjack)
//   - Falls die Summe der Karten 17-20 ist: Stehen
//   - Falls die Summe der Karten 12-16 ist: Stehen, außer wenn der Dealer 7 oder höher hat
//   - Falls die Summe der Karten 11 oder kleiner ist: Karte
func ersterZug(erste string, zweite string, dealer string) string {
	summe := kartenwert(erste) + kartenwert(zweite)
	switch {
	case erste == "Ass" && zweite == "Ass":
		return "Teilen"
	case summe == 21 && kartenwert(dealer) != 10:
		return "Sieg"
	case summe >= 17 && summe <= 20:
		return "Stehen"
	case summe >= 12 && summe <= 16 && kartenwert(dealer) < 7:
		return "Stehen"
	default:
		return "Karte"
	}
}
