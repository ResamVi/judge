package main

import "fmt"

// TODO: Definiere eine 'ofenzeit' Konstante

func main() {
	fmt.Println(ofenzeit)
	fmt.Println(ofenzeit(15))
	fmt.Println(vorbereitungszeit(4))
	fmt.Println(gesamtdauer(4, 8))
}

// TODO: Definiere eine Funktion 'ofenzeit'.
// - Die Funktion akzeptiert einen Parameter, wie lange die Lasagne im Ofen schon ist.
// - Die Funktion gibt eine Zahl als Ergebnis zurück, wie lange die Lasagne noch braucht.

// TODO: Definiere eine Funktion 'vorbereitungszeit'
// - Die Funktion akzeptiert einen Parameter, wie viele Schichten die Lasagne hat
// - Die Funktion gibt zurück, wie lange die Vorbereitungszeit ist abhängig der Anzahl an Schichten
// 		(Vorbereitungszeit = 2 x Anzahl der Schichten)

// TODO: Definiere eine Funktion 'gesamtdauer'
// - Die Funktion akzeptiert zwei Parameter:
//		- Wie lange die Lasagne schon im Ofen ist
//		- Wie viele Schichten die Lasagne hat
// - Die Funktion gibt zurück, wie lange das Kochen der Lasagne bisher gedauert hat
