package main

import (
	"regexp"
	"strings"
	"time"
)

// Contains the date as a string and the time start and end. Example: start_time: 13:00, end_time: 16:00, date: 2022-08-13
type time_slot_struct struct {
	start_time string
	end_time   string
}

// The data from the raflaamo graph api comes with unix timestamps, but we want them as strings.
func convert_unix_timestamp_to_finland(deserialized_graph_data *parsed_graph_data) time_slot_struct {
	time_regex, _ := regexp.Compile(`\d{2}:\d{2}`)

	// Adding 7200000 to the time to match utc +2 (corresponds to 2h)
	unix_start_time_in_finnish_time := time.UnixMilli(int64(deserialized_graph_data.Intervals[0].From + 7200000)).UTC()
	unix_end_time_in_finnish_time := time.UnixMilli(int64(deserialized_graph_data.Intervals[0].To + 7200000)).UTC()

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

// func getCurrentTime() string {
// 	var re, _ = regexp.Compile(`\d{2}:\d{2}`)
// 	dt := time.Now().String()
// 	return re.FindString(dt)
// }
