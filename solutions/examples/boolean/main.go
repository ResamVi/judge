package main

import "fmt"

// istSicher gibt true zurück, wenn die Maschine angefasst werden darf.
// Es ist nicht sicher, wenn:
// - es zu heißt ist
// - Nicht angeschalten ist (Entweder Akku leer oder ausgeschalten)
func istSicher(temperatur int, leer bool, ausgeschalten bool) bool {
	zuHeiß := temperatur > 80

	return !zuHeiß && (leer || ausgeschalten)
}

func main() {
	safe := istSicher(75, false, true)
	fmt.Printf("Sicher: %s\n", safe)
}
