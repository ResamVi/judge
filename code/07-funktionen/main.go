package main

import "fmt"

func PrintHello() {
	fmt.Println("Hallo")
}

// Mit einem Parameter
func PrintHelloName(vorname string, nachname string) {
	fmt.Println("Hallo", vorname, nachname)
}

func main() {
	PrintHelloName("Julien", "Midedji")
}

func Hello(name string) string {
	return "Hallo " + name
}

func HelloAndGoodbye(name string) (string, string) {
	return "Hallo " + name, "Tschüss " + name
}

//func main() {
//	greeting := Hello("Julien")
//	fmt.Println(greeting)
//
//	hello, goodbye := HelloAndGoodbye("Louisa")
//	fmt.Println(hello)
//	fmt.Println(goodbye)
//
//	// Ausgabe:
//	// Hallo Julien
//	// Hallo Louisa
//	// Tschüss Louisa
//}
