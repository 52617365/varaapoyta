package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// error in this function, 87300 returns "0001" when it should return "00:15"
func get_string_time_from_unix(unix_time int64) string {
	time_regex, _ := regexp.Compile(`\d{2}:\d{2}`)

	// Adding 10800000(ms) to the time to match utc +2 or +3 (finnish time) (10800000 ms corresponds to 3h)
	raw_string := time.Unix(unix_time, 0).UTC().String()

	get_time_from_string := time_regex.FindString(raw_string)
	get_time_from_string = strings.Replace(get_time_from_string, ":", "", -1)
	return get_time_from_string
}

func get_unix_from_time(time_to_convert string) int64 {
	// These will not throw an error because we have already validated that they can be converted to integers before.
	minutes, _ := strconv.Atoi(time_to_convert[len(time_to_convert)-2:])
	hour, _ := strconv.Atoi(time_to_convert[:len(time_to_convert)-2])

	// if hour is 0-5 it sets day to 2
	if hour < 5 {
		t := time.Date(1970, time.January, 2, hour, minutes, 00, 0, time.UTC)
		return t.Unix()
	}
	// if hour is 5-23 it sets day to 1
	if hour >= 5 {
		t := time.Date(1970, time.January, 1, hour, minutes, 00, 0, time.UTC)
		return t.Unix()
	}
	return -1
}

func time_formats_are_not_correct(restaurant_starting_time string, restaurant_closing_time string) bool {
	// If time was provided but it was invalid, (can't convert to integer to compare.)
	if len(restaurant_starting_time) > 0 {
		if _, err := strconv.ParseInt(restaurant_starting_time, 10, 64); err != nil {
			return true
		}
	}
	if len(restaurant_closing_time) > 0 {
		if _, err := strconv.ParseInt(restaurant_closing_time, 10, 64); err != nil {
			return true
		}
	}

	// Handle the edge case where the starting_time that was provided is larger than the closing_time.
	if restaurant_starting_time > restaurant_closing_time {
		return true
	}

	// If the times exist but the formats are incorrect meaning their lengths are under 4. Ideal format is something like "1900".
	if restaurant_closing_time != "" && len(restaurant_closing_time) != 4 || restaurant_starting_time != "" && len(restaurant_starting_time) != 4 {
		return true
	}
	return false
}

// Returns all reservation times taking into consideration the restaurants closing time.
// This matters because the restaurants don't take reservations 45 minutes before closing.
// TODO: should restaurant_starting_time and restaurant_closing_time be passed in as unix times?
func get_all_reservation_times(restaurant_starting_time string, restaurant_closing_time string) ([]int64, error) {
	all_times := populate_all_times()

	// if there's no restaurant_starting_time or restaurant_closing_time to take into consideration, just return all the times.
	if restaurant_closing_time == "" && restaurant_starting_time == "" {
		return all_times, nil
	}

	if time_formats_are_not_correct(restaurant_starting_time, restaurant_closing_time) {
		return nil, errors.New("the provided restaurant_starting_time and/or restaurant_closing_time formats were incorrect")
	}

	// if starting time exists but closing does not.
	if opening_time_exists_but_closing_does_not(restaurant_starting_time, restaurant_closing_time) {
		captured_times := extract_unwanted_times(restaurant_starting_time, "", all_times)
		return captured_times, nil
	}
	// if closing time exists but starting time does not.
	if closing_time_exists_but_opening_does_not(restaurant_starting_time, restaurant_closing_time) {
		captured_times := extract_unwanted_times("", restaurant_closing_time, all_times)
		return captured_times, nil
	}
	return nil, errors.New("ok we're here, even though we should not be here")
}

// parameters except all_times are "" if they don't exist.
func extract_unwanted_times(opening_time string, closing_time string, all_times []int64) []int64 {
	if opening_time_exists_but_closing_does_not(opening_time, closing_time) {
		restaurant_starting_time_to_unix := get_unix_from_time(opening_time)
		captured_times := make([]int64, 0, len(all_times))
		for _, individual_time := range all_times {
			// If the time is larger than the restaurants starting time, in other words if the restaurant is already opened
			if restaurant_is_open(individual_time, restaurant_starting_time_to_unix) {
				captured_times = append(captured_times, individual_time)
			}
		}
		return captured_times
	}
	if closing_time_exists_but_opening_does_not(opening_time, closing_time) {
		restaurant_closing_time_to_unix := get_unix_from_time(closing_time)
		captured_times := make([]int64, 0, len(all_times))
		for _, individual_time := range all_times {
			if unix_time_in_not_in_closed_time_slot(restaurant_closing_time_to_unix, individual_time) {
				captured_times = append(captured_times, individual_time)
			}
		}
		return captured_times
	}
	// Here both, opening and closing times exist.
	restaurant_opening_time_to_unix := get_unix_from_time(opening_time)
	restaurant_closing_time_to_unix := get_unix_from_time(closing_time)
	captured_times := make([]int64, 0, len(all_times))
	for _, individual_time := range all_times {
		if unix_time_is_in_between_closing_and_opening_times(restaurant_closing_time_to_unix, individual_time, restaurant_opening_time_to_unix) {
			captured_times = append(captured_times, individual_time)
		}
	}
	return captured_times
}

func unix_time_is_in_between_closing_and_opening_times(restaurant_closing_time_to_unix int64, individual_time int64, restaurant_opening_time_to_unix int64) bool {
	if unix_time_in_not_in_closed_time_slot(restaurant_closing_time_to_unix, individual_time) && individual_time > restaurant_opening_time_to_unix {
		return true
	}
	return false
}

func opening_time_exists_but_closing_does_not(opening_time string, closing_time string) bool {
	if closing_time == "" && opening_time != "" {
		return true
	}
	return false
}

func closing_time_exists_but_opening_does_not(opening_time string, closing_time string) bool {
	if opening_time == "" && closing_time != "" {
		return true
	}
	return false
}

func restaurant_is_open(individual_time int64, restaurant_starting_time_to_unix int64) bool {
	return individual_time > restaurant_starting_time_to_unix
}

// Check to see if the unix timestamp provided is in the time zone where you can't reserve tables (45 minutes before closing) aka 2700 in unix.
func unix_time_in_not_in_closed_time_slot(restaurant_closing_time_to_unix int64, unix_time int64) bool {
	const forty_five_minutes int64 = 2700
	return restaurant_closing_time_to_unix-forty_five_minutes >= unix_time
}

// TODO: This does not return all times correctly E.g. 07:00 is missing. (25200)
func populate_all_times() []int64 {
	all_times := make([]int64, 0, 96)
	minutes := 0
	hour := 0
	for hour < 24 {
		if minutes <= 45 {
			minutes = minutes + 15
		}
		// need to format because number would be "900" instead of "0900".
		if hour < 10 {
			// same thing with minutes
			if minutes < 10 {
				time_to_string := fmt.Sprintf("0%d0%d", hour, minutes)
				all_times = append(all_times, get_unix_from_time(time_to_string))
			}
			if minutes >= 10 {
				time_to_string := fmt.Sprintf("0%d%d", hour, minutes)
				all_times = append(all_times, get_unix_from_time(time_to_string))
			}
		}
		// no need to format.
		if hour >= 10 {
			time_to_string := strconv.Itoa(hour + minutes)
			all_times = append(all_times, get_unix_from_time(time_to_string))
		}
		if minutes == 60 {
			hour++
			minutes = 0
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

// The data from the raflaamo graph api comes as unix timestamps, but we want them as human-readable times in strings, so we
// convert the unix ms timestamps into utc +2 (finnish time).
func convert_unix_timestamp_to_finland_time(time_slot_in_unix *parsed_graph_data) covered_times {
	time_regex, _ := regexp.Compile(`\d{2}:\d{2}`)

	// Adding 10800000(ms) to the time to match utc +2 or +3 (finnish time) (10800000 ms corresponds to 3h)
	unix_start_time_in_finnish_time := time.UnixMilli(int64(time_slot_in_unix.Intervals[0].From + 10800000)).UTC()
	unix_end_time_in_finnish_time := time.UnixMilli(int64(time_slot_in_unix.Intervals[0].To + 10800000)).UTC()

	timestamp_struct_of_available_table := covered_times{
		time:              "",
		time_window_start: strings.Replace(time_regex.FindString(unix_start_time_in_finnish_time.String()), ":", "", -1),
		time_window_end:   strings.Replace(time_regex.FindString(unix_end_time_in_finnish_time.String()), ":", "", -1),
	}

	return timestamp_struct_of_available_table
}

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
// TODO maybe current_time should be unix here?
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
// this throws for some reason.
func time_slots_in_between(start_time string, ending_time string, reservation_times []int64) ([]string, error) {
	if len(start_time) != 4 {
		return nil, errors.New("no start_time in the correct format provided")
	}
	if len(ending_time) != 4 {
		return nil, errors.New("no end_time in the correct format provided")
	}
	start_time = convert_uneven_minutes_to_even(start_time)
	ending_time = convert_uneven_minutes_to_even(ending_time)

	if start_time == "" || ending_time == "" {
		return nil, errors.New("error converting uneven minutes to even minutes")
	}
	if start_time == ending_time {
		return nil, errors.New("trying to get a time_slot with 2 identical timestamps")
	}

	start_time_unix := get_unix_from_time(start_time)
	end_time_unix := get_unix_from_time(ending_time)

	if start_time_unix > end_time_unix {
		return nil, errors.New("start_time_unix was larger than end_time_unix")
	}

	var times_we_want []string
	for _, reservation_time := range reservation_times {
		if reservation_time > start_time_unix && reservation_time <= end_time_unix {
			times_we_want = append(times_we_want, get_string_time_from_unix(reservation_time))
		}
	}

	if len(times_we_want) == 0 {
		return nil, errors.New("there were no times")
	}
	return times_we_want, nil
}
