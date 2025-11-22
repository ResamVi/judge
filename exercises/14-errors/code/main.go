package main

import "fmt"

func main() {
	err := validerePasswort()
	// Prüfe Fehler

	fmt.Println("Passwort angenommen")
}

func validerePasswort() ??? {
	fmt.Println("Gebe ein Passwort ein:")
	var password string
	fmt.Scanln(&password)

	return ???
}
