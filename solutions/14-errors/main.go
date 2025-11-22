package main

import (
	"errors"
	"fmt"
)

func main() {
	err := validerePasswort()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Passwort angenommen")
}

func validerePasswort() error {
	fmt.Println("Gebe ein Passwort ein:")
	var password string
	fmt.Scanln(&password)

	if len(password) <= 8 {
		return errors.New("Passwort ist zu kurz")
	}

	return nil
}
