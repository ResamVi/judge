package main

import "fmt"

func main() {
	var name string = "Max Mustermann"
	var age int = 32
	var height float32 = 1.75

	fmt.Println(name)
	fmt.Println(age)
	fmt.Println(height)

	fmt.Println("Wie ist dein Name?")
	fmt.Scanln(&name)

	fmt.Println("Wie alt bist du?")
	fmt.Scanln(&age)

	fmt.Println("Wie groß bist du?")
	fmt.Scanln(&height)

	fmt.Println(name)
	fmt.Println(age)
	fmt.Println(height)
}
