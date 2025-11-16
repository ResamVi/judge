package main

import "fmt"

func main() {
	fmt.Println(4 + 7)                    // => 11
	fmt.Println(3 - 5)                    // => -2
	fmt.Println(11 * 2)                   // => 22
	fmt.Println(5 / 2)                    // => 2
	fmt.Println(5 % 2)                    // => 1
	fmt.Println(3 < 1)                    // => false
	fmt.Println(7 >= 5)                   // => true
	fmt.Printf("%d + %d = %d", 4, 7, 4+7) // => "4 + 7 = 11"
}
