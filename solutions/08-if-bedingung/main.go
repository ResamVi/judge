package main

import "fmt"

func main() {
	fmt.Println(brauchtFührerschein("auto"))
	fmt.Println(brauchtFührerschein("fahrrad"))
	fmt.Println(brauchtFührerschein("lkw"))

	fmt.Println(schätzeWert(1000, 1))
	fmt.Println(schätzeWert(1000, 5))
	fmt.Println(schätzeWert(1000, 15))
}

// -- Ab hier ändern

// Definiere eine Funktion 'brauchtFührerschein'.
// - Die Funktion nimmt einen Parameter:
//   - Der Name des Fahrzeugs
//
// - Die Funktion soll 'true' zurückgeben für "auto" und "lkw" und sonst 'false'
func brauchtFührerschein(vehicle string) bool {
	if vehicle == "auto" || vehicle == "lkw" {
		return true
	}

	return false
}

// Definiere eine Funktion 'schätzeWert'
// - Die Funktion akzeptiert zwei Parameter:
//   - Den Originalpreis des Autos
//   - Das Alter des Autos
//
// - Wenn das Auto weniger als 3 Jahre alt ist, ist der Schätzwert 80% des Originalpreises
// - Wenn es weniger als 10 Jahre alt ist, ist der Schätzwert 70% des Originalpreises
// - Ist das Auto älter, ist der Schätzwer 50% des Originalpreises
// - Der Rückgabewert soll ein float64 sein
func schätzeWert(originalpreis float32, alter int) float32 {
	if alter < 3 {
		return 0.8 * originalpreis
	}

	if alter < 10 {
		return 0.7 * originalpreis
	}

	return 0.5 * originalpreis
}
