package main

import "fmt"

func main() {
	// TODO: Implement a function TotalBirdCount that accepts a slice of ints that contains the bird count per day.
	// It should return the total number of birds that you counted.
	birdsPerDay := []int{2, 5, 0, 7, 4, 1, 3, 0, 2, 5, 0, 1, 3, 1}
	TotalBirdCount(birdsPerDay)
	// => 34

	birdsPerDay := []int{2, 5, 0, 7, 4, 1, 3, 0, 2, 5, 0, 1, 3, 1}
	BirdsInWeek(birdsPerDay, 2)
	// => 12

	birdsPerDay := []int{2, 5, 0, 7, 4, 1}
	FixBirdCountLog(birdsPerDay)
	// => [3 5 1 7 5 1]
}

// TotalBirdCount return the total bird count by summing
// the individual day's counts.
func TotalBirdCount(birdsPerDay []int) int {
	panic("Please implement the TotalBirdCount() function")
}

// BirdsInWeek returns the total bird count by summing
// only the items belonging to the given week.
func BirdsInWeek(birdsPerDay []int, week int) int {
	panic("Please implement the BirdsInWeek() function")
}

// FixBirdCountLog returns the bird counts after correcting
// the bird counts for alternate days.
func FixBirdCountLog(birdsPerDay []int) []int {
	panic("Please implement the FixBirdCountLog() function")
}
