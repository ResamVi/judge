package main

import "fmt"

func main() {
	tag := 3

	switch tag {
	case 1, 8:
		fmt.Println("Montag")
	case 2:
		fmt.Println("Dienstag")
	case 3:
		fmt.Println("Mittwoch")
	case 4:
		fmt.Println("Donnerstag")
	case 5:
		fmt.Println("Freitag")
	default:
		fmt.Println("Wochenende!")
	}

	punktestand := 72

	if punktestand >= 90 {
		fmt.Println("Note: A")
	} else if punktestand >= 80 {
		fmt.Println("Note: B")
	} else if punktestand >= 70 {
		fmt.Println("Note: C")
	} else if punktestand >= 60 {
		fmt.Println("Note: D")
	} else {
		fmt.Println("Note: F")
	}

	punktestand = 72

	switch {
	case punktestand >= 90:
		fmt.Println("Note: A")
	case punktestand >= 80:
		fmt.Println("Note: B")
	case punktestand >= 70:
		fmt.Println("Note: C")
	case punktestand >= 60:
		fmt.Println("Note: D")
	default:
		fmt.Println("Note: F")
	}

	tag = 6

	switch tag {
	case 1, 2, 3, 4, 5:
		fmt.Println("Es ist Werktags.")
	case 6, 7:
		fmt.Println("Es ist Wochenende!")
	default:
		fmt.Println("Unbekannter Eingabe")
	}

}
