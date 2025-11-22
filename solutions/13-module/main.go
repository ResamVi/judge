package main

import "github.com/sanity-io/litter"

type Person struct {
	Name    string
	Alter   int
	Adresse Adresse
}

type Adresse struct {
	Straße string
	Stadt  string
	PLZ    string
}

func main() {
	litter.Dump(Person{
		Name:  "Max Mustermann",
		Alter: 40,
		Adresse: Adresse{
			Straße: "Mustermannstraße",
			Stadt:  "Karlsruhe",
			PLZ:    "76131",
		},
	})
}
