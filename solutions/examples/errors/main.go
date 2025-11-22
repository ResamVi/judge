package main

import (
	"errors"
	"fmt"
)

func division(a float64, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}

	result := a / b
	return result, nil
}

func main() {
	result, err := division(3, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)

	var number int
	_, err = fmt.Scanln(&number)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(number)
}
