package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// We determine if there is a time slot with open tables by looking at the "color" field in the response.
// The color field will contain "transparent" if it does not contain a graph (open times), else it contains nil (meaning there are open tables)
func time_slot_does_not_contain_open_tables(data *parsed_graph_data) bool {
	return (*data.Intervals)[0].Color == "transparent"
}

func get_opening_and_closing_time_from_kitchen_time(restaurant *response_fields) restaurant_time {
	// Converting restaurant_kitchen_start_time to unix, so we can compare it easily.
	// restaurant_kitchen_start_time := get_unix_from_time(restaurant.Openingtime.Restauranttime.Ranges[0].Start)
	restaurant_kitchen_start_time := get_unix_from_time(restaurant.Openingtime.Kitchentime.Ranges[0].Start)
	// We minus 1 hour from the end time because restaurants don't take reservations before that time slot.
	// IMPORTANT: E.g. if restaurant closes at 22:00, the last possible reservation time is 21:00.
	const one_hour_unix int64 = 3600
	// restaurant_kitchen_ending_time := get_unix_from_time(restaurant.Openingtime.Restauranttime.Ranges[0].End) - one_hour_unix
	restaurant_kitchen_ending_time := get_unix_from_time(restaurant.Openingtime.Kitchentime.Ranges[0].End) - one_hour_unix

	return restaurant_time{
		opening: restaurant_kitchen_start_time,
		closing: restaurant_kitchen_ending_time,
	}
}

func get_string_time_from_unix(unix_time int64) string {
	time_regex, _ := regexp.Compile(`\d{2}:\d{2}`)

	raw_string := time.Unix(unix_time, 0).UTC().String()

	get_time_from_string := time_regex.FindString(raw_string)
	return get_time_from_string
}

func get_unix_from_time(time_to_convert string) int64 {
	time_to_convert = strings.Replace(time_to_convert, ":", "", -1)
	if is_invalid_format(time_to_convert) {
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

// Returns all possible time intervals that can be reserved in the raflaamo reservation page.
// 11:00, 11:15, 11:30 and so on.
func get_all_raflaamo_time_intervals() []int64 {
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

type date_and_time struct {
	date string
	time int64
}

// Gets the current time and date and initializes a struct with it.
func get_current_date_and_time() *date_and_time {
	date_regex, _ := regexp.Compile(`\d{4}-\d{2}-\d{2}`)
	time_regex, _ := regexp.Compile(`\d{2}:\d{2}`)

	dt := time.Now().String()
	date_to_string := date_regex.FindString(dt)
	time_to_string := time_regex.FindString(dt)

	return &date_and_time{
		date: date_to_string,
		time: get_unix_from_time(time_to_string),
	}
}

type covered_times struct {
	time              int64
	time_window_start int64
	time_window_end   int64
}

/*
02:00 covers(00:00-06:00), 08:00 covers(6:00-12:00), 14:00 covers(12:00-18:00), 20:00 covers(18:00-00:00).
The function gets all the time windows we need to check to avoid checking redundant time windows from the past.
*/
func get_graph_time_slots_from_current_point_forward(current_time int64) []covered_times {
	// Getting current_time, so we can avoid checking times from the past.
	all_possible_time_slots := [...]covered_times{
		{time: 7200, time_window_start: 0, time_window_end: 21600},
		{time: 28800, time_window_start: 21600, time_window_end: 43200},
		{time: 50400, time_window_start: 43200, time_window_end: 64800},
		{time: 72000, time_window_start: 64800, time_window_end: 86400},
	}
	time_slots_we_want := make([]covered_times, 0, len(all_possible_time_slots))
	for _, time_slot := range all_possible_time_slots {
		if current_time < time_slot.time_window_end {
			time_slots_we_want = append(time_slots_we_want, time_slot)
		}
	}
	return time_slots_we_want
}
