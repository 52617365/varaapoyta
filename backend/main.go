package main

import "fmt"

// TODO: Make endpoints.
// The line must be drawn somewhere - Jonathan Blow.
func main() {

	current_time := "0000"
	end_time := "0300"
	closing_time := "01:15"

	reservation_times, _ := time_slots_in_between(current_time, end_time, closing_time)

	for _, reservation_time := range reservation_times {
		fmt.Println(reservation_time)
	}
}
