package main

import (
	"regexp"
	"strings"
	"time"
)

// Contains the date as a string and the time start and end. Example: start_time: 13:00, end_time: 16:00, date: 2022-08-13
// @Performance, these could be string references?
type time_slot_struct struct {
	start_time string
	end_time   string
}

// The data from the raflaamo graph api comes as unix timestamps, but we want them as human-readable times in strings, so we
// convert the unix ms timestamps into utc +2 (finnish time).
func convert_unix_timestamp_to_finland_time(time_slot_in_unix *parsed_graph_data) time_slot_struct {
	time_regex, _ := regexp.Compile(`\d{2}:\d{2}`)

	// Adding 7200000(ms) to the time to match utc +2 (finnish time) (7200000 ms corresponds to 2h)
	unix_start_time_in_finnish_time := time.UnixMilli(int64(time_slot_in_unix.Intervals[0].From + 10800000)).UTC()
	unix_end_time_in_finnish_time := time.UnixMilli(int64(time_slot_in_unix.Intervals[0].To + 10800000)).UTC()

	// @Performance, maybe we can get the numbers into the correct format with regex only instead of having to replace ":" with an empty string?
	timestamp_struct_of_available_table := time_slot_struct{
		start_time: strings.Replace(time_regex.FindString(unix_start_time_in_finnish_time.String()), ":", "", -1),
		end_time:   strings.Replace(time_regex.FindString(unix_end_time_in_finnish_time.String()), ":", "", -1),
	}

	return timestamp_struct_of_available_table
}

// 2022-07-21 12:45:54.1414084 +0300 EEST m=+0.001537301
func getCurrentDate() *string {
	re, _ := regexp.Compile(`\d{4}-\d{2}-\d{2}`)
	dt := time.Now().String()
	string_formatted := re.FindString(dt)
	return &string_formatted
}

// 02:00 covers(00:00-06:00), 08:00 covers(6:00-12:00), 14:00 covers(12:00-18:00), 20:00 covers(18:00-00:00)
// The function gets all the time windows we need to check to avoid checking redundant time windows from the past.
func get_time_slots_from_current_point_forward(all_possible_time_slots [4]string) []string {
	current_time, err := get_current_time()
	if err != nil {
		return nil
	}

	// TODO: make sure to cover the timeslots they cover too, -2 and +4 in this loop.
	for index, possible_time_slot := range all_possible_time_slots {
		if possible_time_slot > current_time {
			return all_possible_time_slots[index:]
		}
	}
	return nil
}

/*
Gets current time then calls convert_current_time_to_hours_and_minutes to convert it into a struct which separated hours and
minutes	so that we can work with the time easily.
*/
func get_current_time() (string, error) {
	// TODO: change this regex to not include ":" so we don't have to replace anything ":" with "".
	var re, _ = regexp.Compile(`\d{2}:\d{2}`)
	dt := time.Now().String()
	incorrectly_formatted_time := re.FindString(dt)

	// empty if you can't find a match with regex.
	if incorrectly_formatted_time == "" {
		return "", errors.New("error matching regex in function get_current_hour_and_minutes")
	}

	// Will contain the end result.
	var formatted_time string

	// Add trailing zero if under 1000 because "900" is invalid, we want 0900 instead.
	if len(incorrectly_formatted_time) < 5 {
		formatted_time = fmt.Sprintf("0%s", incorrectly_formatted_time)
	}

	// Reformat E.g. 10:00 to 1000.
	formatted_time = strings.Replace(incorrectly_formatted_time, ":", "", -1)

	return formatted_time, nil
}
