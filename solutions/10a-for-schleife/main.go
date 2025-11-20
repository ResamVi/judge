package main

import "fmt"

func main() {
	fmt.Println("Nenne eine Zahl:")

	var n int
	fmt.Scanln(&n)

	for i := 1; i <= n; i++ {
		fmt.Println(i)
	}
}
