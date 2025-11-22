package main

import "fmt"

type Film struct {
	Titel string
	Jahr  int
	Genre string
}

type Rechteck struct {
	Höhe   int
	Breite int
}

type Konto struct {
	Name     string
	Guthaben int
}

func main() {
	fmt.Println(flächeBerechnen(Rechteck{Höhe: 12, Breite: 5}))
	fmt.Println(werbungZeigen(Film{
		Titel: "Matrix",
		Jahr:  1999,
		Genre: "ScienceFiction",
	}))
	fmt.Printf("%#v\n", kontoÖffnen("Svenja Schmidt", 100))
}

func flächeBerechnen(rechteck Rechteck) int {
	return rechteck.Höhe * rechteck.Breite
}

func werbungZeigen(film Film) string {
	return fmt.Sprintf("Der Film %v kam im Jahr %v und gehört zur %v Genre", film.Titel, film.Jahr, film.Genre)
}

func kontoÖffnen(name string, guthaben int) Konto {
	return Konto{Name: name, Guthaben: guthaben}
}
