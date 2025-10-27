package main

import "fmt"

func main() {
	fmt.Println(kartenwert("ass"))
	fmt.Println(kartenwert("zehn"))
	fmt.Println(kartenwert("drei"))
	fmt.Println(kartenwert("könig"))
	fmt.Println(kartenwert("bube"))
	fmt.Println(kartenwert("dame"))
}

// TODO: Definiere eine Funktion 'kartenwert'.
// - Die Funktion nimmt einen Parameter:
//		- Der Name der Karte
// - Die Funktion soll den Wert der Karte zurückgeben

// TODO: Definiere eine Funktion 'ersterZug'
// - Die Funktion nimmt drei Parameter:
//		- Der Name deiner ersten Karte
//		- Der Name deiner zweiten Karte
//		- Der Name der Karte des Dealers, die offen ist
// - Die Funktion gibt den String zurück welcher Zug gemacht werden soll
//		- "Karte" (eine weitere Karte nehmen)
//		- "Stehen" (keine weitere mehr nehmen)
//		- "Teilen" (Karten aufteilen)
//		- "Sieg" (wenn man schon 21 hat)
// - Die Funktion soll folgende Strategie implementieren:
//		- Falls zwei Asse: Teilen
//		- Falls dein Kartenwert 21 ist und der Dealer kein Ass/10/Bube/Dame/König: Sieg (weil Blackjack)
//		- Falls die Summe der Karten 17-20 ist: Stehen
//		- Falls die Summe der Karten 12-16 ist: Stehen, außer wenn der Dealer 7 oder höher hat
//		- Falls die Summe der Karten 11 oder kleiner ist: Karte
//
