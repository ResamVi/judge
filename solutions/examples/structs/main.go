package main

import "fmt"

type Person struct {
	Name  string
	Alter int
}

func main() {
	markus := Person{Name: "Markus", Alter: 42}
	fmt.Println(markus) // => {Markus 42}

	markus.Name = "Markus Schmidt"
	fmt.Printf("Name: %s | Alter: %d\n", markus.Name, markus.Alter)
	// Output: Name: Markus Schmidt | Alter: 42

	markus = altern(markus)
	fmt.Printf("Name: %s | Alter: %d\n", markus.Name, markus.Alter)
	// Output: Name: Markus Schmidt | Alter: 43
}

func altern(mensch Person) Person {
	mensch.Alter += 1
	return mensch
}
