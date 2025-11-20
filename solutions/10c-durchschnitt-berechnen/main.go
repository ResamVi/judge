package main

import "fmt"

func main() {
	fmt.Println("Gebe mir Zahlen:")

	sum := 0.0
	count := 0.0
	for {
		var eingabe float64
		fmt.Scanln(&eingabe)

		if eingabe == 0 {
			break
		}

		if eingabe < 0 {
			continue
		}

		sum += eingabe
		count++
	}

	fmt.Println(sum / count)
}
