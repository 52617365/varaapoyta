package main

import (
	"fmt"
	"regexp"
	"time"
)

// Contains the date as a string and the time start and end. Example: start_time: 13:00, end_time: 16:00, date: 2022-08-13
type time_slot_struct struct {
	start_time string
	end_time   string
}

// FIX: unix time conversion does not yet work correctly; fix it.
// The data from the raflaamo graph api comes with unix timestamps, but we want them as strings.
func convert_unix_timestamp_to_finland(deserialized_graph_data *parsed_graph_data) time_slot_struct {
	// TODO: convert unix timestamp to finnish timezone stamp. (utc+2)
	unix_start_time_in_finnish_time := time.UnixMilli(int64(deserialized_graph_data.Intervals[0].From + 7200000)).UTC()
	unix_end_time_in_finnish_time := time.UnixMilli(int64(deserialized_graph_data.Intervals[0].To + 7200000)).UTC()

	fmt.Println(unix_start_time_in_finnish_time)
	fmt.Println(unix_end_time_in_finnish_time)
	timestamp_struct_of_available_table := time_slot_struct{
		start_time: unix_start_time_in_finnish_time.String(),
		end_time:   unix_end_time_in_finnish_time.String(),
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

func getCurrentTime() string {
	var re, _ = regexp.Compile(`\d{2}:\d{2}`)
	dt := time.Now().String()
	return re.FindString(dt)
}
