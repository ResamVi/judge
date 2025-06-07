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
// - Die Funktion gibt einen Buchstaben zurück für den Zug der gemacht werden soll
//		- "Karte" (eine weitere Karte nehmen)
//		- "Stehen" (keine weitere mehr nehmen)
//		- "Teilen" (Karten aufteilen)
//		- "Sieg" (wenn man schon 21 hat)
// - Die Funktion soll folgende Strategie implementieren
//		- 
//
