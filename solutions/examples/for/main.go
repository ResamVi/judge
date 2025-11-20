package main

import "fmt"

func main() {
	for i := 0; i <= 5; i++ {
		fmt.Println(i)
	}

	x := 0
	for x < 5 {
		fmt.Println(x)
		x++
	}

	a := 10
	fmt.Println(a) // 10
	a++            // gleich wie: a += 1 oder a = a + 1
	fmt.Println(a) // 11

	for {
		fmt.Println("Bis in die Ewigkeit...")
	}

	for {
		var eingabe int
		fmt.Scanln(&eingabe)

		if eingabe < 0 {
			fmt.Println("Keine negativen Eingaben erlaubt!")
			continue
		}

		if eingabe == 0 {
			break
		}

		fmt.Println(eingabe)
	}

	//for <Variable initialisieren>; <Bedingung>; <Variable ändern> {
	//	// Code der wiederholt wird solange <Bedingung> wahr ist
	//}
}
