package main

import (
	"fmt"
)

func main() {
	lines := [][]rune{
		{32, 47, 92, 95, 47, 92},
		{40, 32, 111, 46, 111, 32, 41},
		{32, 62, 32, 94, 32, 60},
	}

	for _, line := range lines {
		fmt.Println(string(line))
	}
}
