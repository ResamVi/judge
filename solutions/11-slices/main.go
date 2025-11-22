package main

import "fmt"

func main() {
	fmt.Println(stapel())
	fmt.Println(karteNehmen(stapel(), 3))
	fmt.Println(karteAustauschen(stapel(), 13, "K"))
	fmt.Println(karteHinzufügen(stapel(), "1"))
}

func stapel() []string {
	return []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
}

func karteNehmen(stapel []string, karte int) string {
	return stapel[karte]
}

func karteAustauschen(stapel []string, position int, karte string) []string {
	stapel[position] = karte
	return stapel
}

func karteHinzufügen(stapel []string, karte string) []string {
	stapel = append(stapel, karte)
	return stapel
}
