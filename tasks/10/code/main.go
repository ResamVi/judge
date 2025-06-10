package main

import "fmt"

func main() {
	fmt.Println(geradeKarten()) 
	fmt.Println(zieheKarten([]int{1,2,3}, 2)) 
	fmt.Println(zieheKarten([]int{1,2,3,4,5}, 12)) 
	fmt.Println(ersetzeKarte([]int{1,2,3,4,5}, 1, 10)) 
	fmt.Println(ersetzeKarte([]int{1,2,3,4,5}, 12, 10)) 
}

// TODO: Definiere eine Funktion "geradeKarten"
// - Die Funktion gibt ein slice zurück von allen Karten die gerade sind

// TODO: Definiere eine Funktion 'zieheKarte'. 
// - Die Funktion nimmt zwei Parameter:
//		- Ein Stapel an Karten ([]int)
//		- Die Position der Karte die aus dem Stapel gezogen wird (int)
// - Die Funktion gibt die Karte zurück die an der übergebenenen Position ist
// - Wenn die Zahl größer ist als die Anzahl der Karten im Slice gib -1 zurück

// TODO: Definiere eine Funktion 'ersetzeKarte'
// - Die Funktion nimmt drei Parameter:
//		- Ein Stapel an Karten ([]int)
//		- Die Position der Karte die ersetzt wird (int)
//		- Die neue Karte mit der sie ersetzt wird (int)
// - Die Funktion ersetzt an der Position des Slices die Karte mit der neuen Zahl und gibt das neue Slice zurück

// TODO: Definiere eine Funktion 'ergänzeKarte'
// - Die Funktion nimmt zwei Parameter:
//		- Ein Stapel an Karten ([]int)
//		- Die neue Karte die ergänzt wird
// - Die Funktion gibt den Stapel zurück mit der neuen Karte an der letzten Position
