package main

import "fmt"

func main() {
	fmt.Println("Was ist deine Note in Geometrie?")
	var geometrie float32
	fmt.Scanln(&geometrie)

	fmt.Println("Was ist deine Note in Algebra?")
	var algebra float32
	fmt.Scanln(&algebra)

	fmt.Println("Was ist deine Note in Physik?")
	var physik float32
	fmt.Scanln(&physik)

	durchschnitt := (geometrie + algebra + physik) / 3.0

	fmt.Println("Dein Notendurchschnitt:")
	fmt.Println(durchschnitt)
	fmt.Println(durchschnitt <= 4)
}
