package main

import (
	"strconv"
	"strings"
)

type time_utils struct {
	current_time date_and_time
	closing_time int64
}

type relative_time struct {
	hour    int
	minutes int
}

func (t time_utils) get_time_till_restaurant_closing_time() relative_time {
	// we minus one hour from it cuz they don't take reservations in that time slot, but they're still technically open, so we add it back here, this is the only place where we add it back.
	const one_hour_unix int64 = 3600
	closing_time := t.closing_time
	current_time := t.current_time

	closing_time += one_hour_unix

	// If the following condition hits, restaurant is already closed.
	if closing_time <= current_time.time {
		return relative_time{hour: -1, minutes: -1}
	}

	// TODO: Experimental fix already in place, see if it works, if not, revisit. When closing time is 81000 (23:30) and current_time.time is 75000, the subtraction should be 6000 and not 43500 that is currently is.
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

func (t time_utils) get_time_left_to_reserve() relative_time {
	closing_time := t.closing_time
	current_time := t.current_time
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
