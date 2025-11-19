package main

import "fmt"

func main() {
	personen := 20
	katzen := 30
	hunde := 15

	if personen < katzen {
		fmt.Println("Zu viele Katzen!")
	}

	hunde += 5

	if personen == hunde {
		fmt.Println("Gleich viele Menschen wie Hunde")
	}

	result := "Die Zahl ist... "

	var number int
	if number > 0 {
		result += "positiv"
	} else if number < 0 {
		result += "negativ"
	} else {
		result += "null"
	}
	fmt.Println(result)

}
