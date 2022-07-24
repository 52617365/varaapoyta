package main

import "fmt"

// TODO: get data first and then make the endpoints.

func main() {
	times := getAllPossibleTimes()
	//getAvailableTables()
	fmt.Println(times)
}
