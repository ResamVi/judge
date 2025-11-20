package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println("Versuche die Zahl zu erraten:")

	zahl := rand.Intn(10)
	versuche := 0
	for {
		var eingabe int
		fmt.Scanln(&eingabe)

		if eingabe == zahl {
			break
		}
		versuche++
	}

	fmt.Printf("Richtig! Anzahl der Versuche: %v\n", versuche)
}
