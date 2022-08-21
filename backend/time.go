package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func get_string_time_from_unix(unix_time int64) string {
	time_regex, _ := regexp.Compile(`\d{2}:\d{2}`)

	raw_string := time.Unix(unix_time, 0).UTC().String()

	get_time_from_string := time_regex.FindString(raw_string)
	return get_time_from_string
}

func get_unix_from_time(time_to_convert string) int64 {
	time_to_convert = strings.Replace(time_to_convert, ":", "", -1)
	if is_not_valid_format(time_to_convert) {
		return -1
	}

	minutes, _ := strconv.Atoi(time_to_convert[len(time_to_convert)-2:])
	hour, _ := strconv.Atoi(time_to_convert[:len(time_to_convert)-2])

	// if hour is 0-5 it sets day to 2 (unix)
	if hour < 5 {
		t := time.Date(1970, time.January, 2, hour, minutes, 00, 0, time.UTC)
		return t.Unix()
	}
	// if hour is 5-23 it sets day to 1 (unix)
	if hour >= 5 {
		t := time.Date(1970, time.January, 1, hour, minutes, 00, 0, time.UTC)
		return t.Unix()
	}
	return -1
}

// The parameters passed here have already taken into consideration the restaurants starting and closing time.
func get_all_reservation_times(restaurant_starting_time_unix int64, restaurant_closing_time_unix int64) ([]int64, error) {
	all_times := populate_all_times()

	// Here we check if the starting_time is larger than closing_time.
	// This will result in an error because the user tried to provide times but failed with the format.
	if restaurant_starting_time_unix >= restaurant_closing_time_unix {
		return nil, errors.New("restaurant_start_time was larger or equal to closing_time")
	}
	captured_times := extract_unwanted_times(restaurant_starting_time_unix, restaurant_closing_time_unix, all_times)
	return captured_times, nil
}

func extract_unwanted_times(first_possible_reservation_time int64, last_possible_reservation_time int64, all_times []int64) []int64 {
	captured_times := make([]int64, 0, len(all_times))

	for _, time := range all_times {
		if time > first_possible_reservation_time && time <= last_possible_reservation_time {
			captured_times = append(captured_times, time)
		}
	}
	return captured_times
}

func populate_all_times() []int64 {
	all_times := make([]int64, 0, 96)
	for hour := 0; hour < 24; hour++ {
		if hour < 10 {
			formatted_time_slot_one := get_unix_from_time(fmt.Sprintf("0%d00", hour))
			formatted_time_slot_two := get_unix_from_time(fmt.Sprintf("0%d15", hour))
			formatted_time_slot_three := get_unix_from_time(fmt.Sprintf("0%d30", hour))
			formatted_time_slot_four := get_unix_from_time(fmt.Sprintf("0%d45", hour))
			all_times = append(all_times, formatted_time_slot_one)
			all_times = append(all_times, formatted_time_slot_two)
			all_times = append(all_times, formatted_time_slot_three)
			all_times = append(all_times, formatted_time_slot_four)
		}
		if hour >= 10 {
			formatted_time_slot_one := get_unix_from_time(fmt.Sprintf("%d00", hour))
			formatted_time_slot_two := get_unix_from_time(fmt.Sprintf("%d15", hour))
			formatted_time_slot_three := get_unix_from_time(fmt.Sprintf("%d30", hour))
			formatted_time_slot_four := get_unix_from_time(fmt.Sprintf("%d45", hour))
			all_times = append(all_times, formatted_time_slot_one)
			all_times = append(all_times, formatted_time_slot_two)
			all_times = append(all_times, formatted_time_slot_three)
			all_times = append(all_times, formatted_time_slot_four)
		}
	}
	return all_times
}

// This struct contains the time you check the graph api with, and the corresponding start and end time window that the response covers.
type covered_times struct {
	time              string
	time_window_start string
	time_window_end   string
}

//// The data from the raflaamo graph api comes as unix timestamps, but we want them as human-readable times in strings, so we
//// convert the unix ms timestamps into utc +2 (finnish time).
//func convert_unix_timestamp_to_finland_time(time_slot_in_unix *parsed_graph_data) covered_times {
//	time_regex, _ := regexp.Compile(`\d{2}:\d{2}`)
//
//	// Adding 10800000(ms) to the time to match utc +2 or +3 (finnish time) (10800000 ms corresponds to 3h)
//	unix_start_time_in_finnish_time := time.UnixMilli(int64(time_slot_in_unix.Intervals[0].From + 10800000)).UTC()
//	unix_end_time_in_finnish_time := time.UnixMilli(int64(time_slot_in_unix.Intervals[0].To + 10800000)).UTC()
//
//	timestamp_struct_of_available_table := covered_times{
//		time:              "",
//		time_window_start: strings.Replace(time_regex.FindString(unix_start_time_in_finnish_time.String()), ":", "", -1),
//		time_window_end:   strings.Replace(time_regex.FindString(unix_end_time_in_finnish_time.String()), ":", "", -1),
//	}
//
//	return timestamp_struct_of_available_table
//}

type date_and_time struct {
	date string
	time string
}

// 2022-07-21 12:45:54.1414084 +0300 EEST m=+0.001537301
func get_current_date_and_time() date_and_time {
	date_regex, _ := regexp.Compile(`\d{4}-\d{2}-\d{2}`)
	time_regex, _ := regexp.Compile(`\d{2}:\d{2}`)

	dt := time.Now().String()
	date_to_string := date_regex.FindString(dt)
	time_to_string := time_regex.FindString(dt)
	time_to_string = strings.Replace(time_to_string, ":", "", -1) // Reformats E.g. 10:00 to 1000.

	return date_and_time{
		date: date_to_string,
		time: time_to_string,
	}
}

// 02:00 covers(00:00-06:00), 08:00 covers(6:00-12:00), 14:00 covers(12:00-18:00), 20:00 covers(18:00-00:00).
// The function gets all the time windows we need to check to avoid checking redundant time windows from the past.
func get_time_slots_from_current_point_forward(current_time string) []covered_times {
	all_possible_time_slots := [...]covered_times{
		{time: "0200", time_window_start: "0000", time_window_end: "0600"},
		{time: "0800", time_window_start: "0600", time_window_end: "1200"},
		{time: "1400", time_window_start: "1200", time_window_end: "1800"},
		{time: "2000", time_window_start: "1800", time_window_end: "0000"},
	}
	time_slots_we_want := make([]covered_times, 0, len(all_possible_time_slots))
	for _, time_slot := range all_possible_time_slots {
		if current_time < time_slot.time_window_end {
			time_slots_we_want = append(time_slots_we_want, time_slot)
		}
	}
	return time_slots_we_want
}

/*
Used to get all the time slots in between the graph start and graph end.
E.g. if start is 2348 and end is 0100, it will get time slots 0000, 0015, 0030, 0045, 0100.
*/
// Here reservation_times here has already taken into consideration the restaurants opening and closing time.
func time_slots_in_between(start_time int64, graph_end int64, reservation_times []int64) ([]string, error) {
	//if start_time == -1 || graph_end == -1 {
	//	return nil, errors.New("start_time or graph_end were empty")
	//}
	if start_time == graph_end {
		return nil, errors.New("trying to get a time_slot with 2 identical timestamps")
	}
	if start_time > graph_end {
		return nil, errors.New("start_time was larger than end_time")
	}

	var times_we_want []string
	for _, reservation_time := range reservation_times {
		if reservation_time > start_time && reservation_time <= graph_end {
			// We convert the times into string_time because that's the format we will be using later on to display the times.
			times_we_want = append(times_we_want, get_string_time_from_unix(reservation_time))
		}
	}
	if len(times_we_want) == 0 {
		return nil, errors.New("no times found")
	}
	return times_we_want, nil
}
