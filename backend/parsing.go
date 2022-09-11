package main

import (
	"strconv"
	"strings"
)

type relative_time struct {
	hour    int
	minutes int
}

func get_time_till_restaurant_closing_time(closing_time int64) relative_time {
	// we minused one hour from it cuz they don't take reservations in that time slot, but they're still technically open, so we add it back here, this is the only place where we add it back.
	const one_hour_unix int64 = 3600
	closing_time += one_hour_unix
	current_time := get_current_date_and_time()
	// already closed.
	if closing_time <= current_time.time {
		return relative_time{hour: -1, minutes: -1}
	}

	time_left_to_closing_unix := closing_time - current_time.time

	relative_time_string := get_string_time_from_unix(time_left_to_closing_unix)
	relative_time_string = strings.Replace(relative_time_string, ":", "", -1)

	if is_not_valid_format(relative_time_string) {
		return relative_time{hour: -1, minutes: -1}
	}

	minutes, _ := strconv.Atoi(relative_time_string[len(relative_time_string)-2:])
	hour, _ := strconv.Atoi(relative_time_string[:len(relative_time_string)-2])

	return relative_time{
		hour:    hour,
		minutes: minutes,
	}

}
func is_not_valid_format(our_number string) bool {
	if _, err := strconv.ParseInt(our_number, 10, 64); err != nil {
		return true
	}
	if len(our_number) != 4 {
		return true
	}
	if our_number == "" {
		return true
	}
	return false
}
