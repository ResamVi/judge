package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	for zeile := 1; zeile <= 5; zeile++ {
		ausgabe(zeile)
		time.Sleep(100 * time.Millisecond)
	}
	for zeile := 5; zeile >= 1; zeile-- {
		ausgabe(zeile)
		time.Sleep(100 * time.Millisecond)
	}
}

func ausgabe(laenge int) {
	ergebnis := strings.Repeat("#", laenge)
	fmt.Println(ergebnis)
}
