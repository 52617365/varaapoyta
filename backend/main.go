package main

import "fmt"

// TODO: Make endpoints.

// Fuck times stored in strings.
// The line must be drawn here - Jonathan Blow.
func main() {
	current_time := "0000"
	end_time := "0300"
	closing_time := "01:15"

	all_reservation_times := get_all_reservation_times(closing_time) // in reality, it's not all because we need to consider restaurants closing time.
	reservation_times, _ := time_slots_in_between(current_time, end_time, all_reservation_times)

	for _, reservation_time := range reservation_times {
		fmt.Println(reservation_time)
	}
}
