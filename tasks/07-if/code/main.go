package main

import "fmt"

func main() {
	fmt.Println(brauchtFührerschein("auto"))
	fmt.Println(brauchtFührerschein("fahrrad"))
	fmt.Println(brauchtFührerschein("lkw"))

	fmt.Println(wähleAuto("Aston Martin Valhalla", "Tesla Model 3"))
	fmt.Println(wähleAuto("Honda Civic", "Ferrari F80"))
	fmt.Println(wähleAuto("ford", "Bugatti"))

	fmt.Println(schätzeWert(1000, 1))
	fmt.Println(schätzeWert(1000, 5))
	fmt.Println(schätzeWert(1000, 15))
}

// TODO: Definiere eine Funktion 'brauchtFührerschein'.
// - Die Funktion nimmt einen Parameter:
//		- Der Name des Fahrzeugs
// - Die Funktion soll 'true' zurückgeben für "auto" und "lkw" und sonst 'false'

// TODO: Definiere eine Funktion 'wähleAuto'
// - Die Funktion nimmt zwei Parameter:
//		- Den Namen des ersten Autos
//		- Den Namen des zweiten Autos
// - Die Funktion gibt den Namen des Autos zurück, welches lexikographisch zuerst kommt
//	 (Audi kommt vor BMW, weil A < B)
//

// TODO: Definiere eine Funktion 'schätzeWert'
// - Die Funktion akzeptiert zwei Parameter:
//		- Den Originalpreis des Autos
//		- Das Alter des Autos
// - Wenn das Auto weniger als 3 Jahre alt ist, ist der Schätzwert 80% des Originalpreises
// - Wenn es weniger als 10 Jahre alt ist, ist der Schätzwert 70% des Originalpreises
// - Ist das Auto älter, ist der Schätzwer 50% des Originalpreises
// - Der Rückgabewert soll ein float64 sein
