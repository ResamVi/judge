package main

import "fmt"

func main() {
	fmt.Println("Nenne eine Zahl:")

	var n int
	fmt.Scanln(&n)

	ergebnis := 1
	for i := 1; i <= n; i++ {
		ergebnis *= i
	}
	fmt.Println(ergebnis)
}
