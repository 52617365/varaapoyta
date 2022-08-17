package main

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type time_slot_window struct {
	time              string
	time_window_start string
	time_window_end   string
}

// 02:00 covers(00:00-06:00), 08:00 covers(6:00-12:00), 14:00 covers(12:00-18:00), 20:00 covers(18:00-00:00)
func get_all_time_windows(current_time string) []time_slot_window {
	time_windows := [...]time_slot_window{
		{time: "0200", time_window_start: "0000", time_window_end: "0600"},
		{time: "0800", time_window_start: "0600", time_window_end: "1200"},
		{time: "1400", time_window_start: "1200", time_window_end: "1800"},
		{time: "2000", time_window_start: "1800", time_window_end: "0000"},
	}
	time_windows_from_current_forward := get_time_slots_from_current_point_forward(time_windows, current_time)
	return time_windows_from_current_forward
}

// Contains the date as a string and the time start and end. Example: start_time: 13:00, end_time: 16:00
type graph_time_slot struct {
	start_time string
	end_time   string
}

// The data from the raflaamo graph api comes as unix timestamps, but we want them as human-readable times in strings, so we
// convert the unix ms timestamps into utc +2 (finnish time).
func convert_unix_timestamp_to_finland_time(time_slot_in_unix *parsed_graph_data) graph_time_slot {
	time_regex, _ := regexp.Compile(`\d{2}:\d{2}`)

	// Adding 7200000(ms) to the time to match utc +2 (finnish time) (7200000 ms corresponds to 2h)
	unix_start_time_in_finnish_time := time.UnixMilli(int64(time_slot_in_unix.Intervals[0].From + 10800000)).UTC()
	unix_end_time_in_finnish_time := time.UnixMilli(int64(time_slot_in_unix.Intervals[0].To + 10800000)).UTC()

	// @Performance, maybe we can get the numbers into the correct format with regex only instead of having to replace ":" with an empty string?
	timestamp_struct_of_available_table := graph_time_slot{
		start_time: strings.Replace(time_regex.FindString(unix_start_time_in_finnish_time.String()), ":", "", -1),
		end_time:   strings.Replace(time_regex.FindString(unix_end_time_in_finnish_time.String()), ":", "", -1),
	}

	return timestamp_struct_of_available_table
}

// 2022-07-21 12:45:54.1414084 +0300 EEST m=+0.001537301
func get_current_date() string {
	re, _ := regexp.Compile(`\d{4}-\d{2}-\d{2}`)
	dt := time.Now().String()
	string_formatted := re.FindString(dt)
	return string_formatted
}

// 02:00 covers(00:00-06:00), 08:00 covers(6:00-12:00), 14:00 covers(12:00-18:00), 20:00 covers(18:00-00:00).
// The function gets all the time windows we need to check to avoid checking redundant time windows from the past.
func get_time_slots_from_current_point_forward(all_possible_time_slots [4]time_slot_window, current_time string) []time_slot_window {
	for time_slot_index, time_slot := range all_possible_time_slots {
		if current_time < time_slot.time_window_end {
			return all_possible_time_slots[time_slot_index:]
		}
	}
	return nil
}

/*
Gets current time then calls convert_current_time_to_hours_and_minutes to convert it into a struct which separated hours and
minutes	so that we can work with the time easily.
*/
func get_current_time() (string, error) {
	var re, _ = regexp.Compile(`\d{2}:\d{2}`)
	dt := time.Now().String()
	incorrectly_formatted_time := re.FindString(dt)

	// empty if you can't find a match with regex.
	if incorrectly_formatted_time == "" {
		return "", errors.New("error matching regex in function get_current_hour_and_minutes")
	}

	// Reformat E.g. 10:00 to 1000.
	formatted_time := strings.Replace(incorrectly_formatted_time, ":", "", -1)

	return formatted_time, nil
}
