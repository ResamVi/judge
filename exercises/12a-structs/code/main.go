package main

import "fmt"

type Film struct {
	Titel string
	Jahr int
	Genre string
}

func main() {
	fmt.Println(flächeBerechnen(Rechteck{Höhe: 12, Breite: 5}))
	fmt.Println(werbungZeigen(Film{ ??? })) // Erwartete Ausgabe: "Der Film <FILM_NAME_HIER> kam im Jahr <JAHR_HIER> und gehört zur <GENRE_HIER> Genre"
	fmt.Printf("%#v\n", kontoÖffnen("Svenja Schmidt", 100))
}

func flächeBerechnen(???) ??? {

}

func werbungZeigen(film Film) string {
	return ""
}

func kontoÖffnen(name string, guthaben int) ??? {
	return ???
}