package main

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
	// Erstelle eine Person "Max Mustermann" mit 40 Jahren
	// welcher in der Mustermannstraße in Karlsruhe 76131 wohnt
	//
	// Gebe dann das struct mit dieser Library aus:
	// https://github.com/sanity-io/litter
	// benutze dafür litter.Dump(...)

}
