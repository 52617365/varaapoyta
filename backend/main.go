package main

import (
	"fmt"
)

// TODO: Make endpoints.

// Fuck times stored in strings.
// The line must be drawn here - Jonathan Blow.
func main() {
	current_time := "2330"
	end_time := "0259"
	closing_time := "01:30"
	starting_time := "11:00"

	all_reservation_times, _ := get_all_reservation_times(starting_time, closing_time)
	reservation_times, _ := time_slots_in_between(current_time, end_time, all_reservation_times)

	for _, reservation_time := range reservation_times {
		fmt.Println(reservation_time)
	}
}
