package main

import (
	"fmt"
)

// TODO: Make endpoints.

// Fuck times stored in strings.
// The line must be drawn here - Jonathan Blow.
func main() {
	current_time := "2359"
	end_time := "0259"
	closing_time := "01:30"

	// all reservation_times works as intended.
	all_reservation_times := get_all_reservation_times(closing_time)
	// this does not.
	reservation_times, _ := time_slots_in_between(current_time, end_time, all_reservation_times)

	for _, reservation_time := range reservation_times {
		fmt.Println(reservation_time)
	}
}
