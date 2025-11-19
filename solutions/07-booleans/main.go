package main

import "fmt"

func main() {
	fmt.Printf("Hat ein Ticket: %v\n", eintrittErlaubt(true, false))
	fmt.Printf("Ist VIP: %v\n", eintrittErlaubt(false, true))
	fmt.Printf("Hat weder Ticket noch VIP: %v\n", eintrittErlaubt(false, false))

	fmt.Printf("Nicht eingesteckt: %v\n", computerLäuft(false, true))
	fmt.Printf("Nicht angeschalten: %v\n", computerLäuft(true, false))
	fmt.Printf("Eingesteckt und angeschalten: %v\n", computerLäuft(true, true))

	fmt.Printf("Name enthält Zahlen: %v\n", nameValide(true))
	fmt.Printf("Name enthält keine Zahlen: %v\n", nameValide(false))

	fmt.Printf("101 Grad: %v\n", istHeiß(101))
	fmt.Printf("100 Grad: %v\n", istHeiß(100))
	fmt.Printf("99 Grad: %v\n", istHeiß(99))
}

// Erste Funktion: eintrittErlaubt(hatTicket bool, istVIP bool) bool
// Ausgabe: Der Eintritt ist erlaubt, wenn die Person ein Ticket hat oder VIP ist
func eintrittErlaubt(hatTicket bool, istVIP bool) bool {
	return hatTicket || istVIP
}

// Zweite Funktion: computerLäuft(eingesteckt bool, angeschalten bool) bool
// Ausgabe: Der Computer läuft, wenn der Stecker eingesteckt ist und angeschalten wurde
func computerLäuft(eingesteckt bool, angeschalten bool) bool {
	return eingesteckt && angeschalten
}

// Dritte Funktion: nameValide(enthältZahlen bool) bool
// Ausgabe: Ein Name ist valide, wenn es keine Zahlen beinhaltet
func nameValide(enthähltZahlen bool) bool {
	return !enthähltZahlen
}

// Vierte Funktion: istHeiß(temperatur int) bool
// Ausgabe: Die Temperatur ist heiß wenn sie über 100 Grad enthält
func istHeiß(temperatur int) bool {
	return temperatur > 100
}
