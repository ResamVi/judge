package main

import "fmt"

func main() {
	var meinSlice []int // Leer
	fmt.Println(meinSlice)
	//meinSlice := []int{9, 10, 11, 12, 13, 14} // mit Daten gefüllt

	fmt.Println(meinSlice[3]) // 12
	fmt.Println(meinSlice[0]) // 9
	fmt.Println(meinSlice[6]) // panic: runtime error: index out of range [6] with length 6

	meinSlice[1] = 5
	fmt.Println(meinSlice)

	a := []int{1, 3}
	a = append(a, 4)
	// => []int{1,3,4}

	//len([]string{"hallo", "wie", "gehts"}) // 3

	for i := 0; i < len(a); i++ {
		fmt.Println(a[i] + 1)
	}

}
