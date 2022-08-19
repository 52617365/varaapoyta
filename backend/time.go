package main

import (
	"errors"
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

func get_unix_from_time(hour int, minutes int) int64 {
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

// Returns all reservation times taking into consideration the restaurants closing time.
// This matters because the restaurants don't take reservations 45 minutes before closing.

// TODO: handle opening time too on top of closing_time. (Don't take times from before the opening time.)
// TODO: fix problem where restaurant_closing_time set to 0115 returns times until 2330 when it should be 0030.
func get_all_reservation_times(restaurant_closing_time string) []int64 {
	all_times := populate_all_times()

	// last 2 letters.
	mins := restaurant_closing_time[len(restaurant_closing_time)-2:]
	hours := restaurant_closing_time[:len(restaurant_closing_time)-2]

	mins_int, err := strconv.Atoi(mins)
	if err != nil {
		return nil
	}
	hour_int, err := strconv.Atoi(hours)
	if err != nil {
		return nil
	}

	// 0000 to 0800 could be 2 jan Mon, 2 Jan 1970 00:00:00 GMT
	// 0800 to 2345 could be 1 jan Mon, 1 Jan 1970 00:00:00 GMT
	// store as unix, return as string?

	restaurant_closing_time_to_unix := get_unix_from_time(hour_int, mins_int)

	captured_times := make([]int64, 0, len(all_times))
	for _, individual_time := range all_times {
		if unix_time_in_not_in_closing_time_slot(restaurant_closing_time_to_unix, individual_time) {
			captured_times = append(captured_times, individual_time)
		}
	}
	return captured_times
}

// Check to see if the unix timestamp provided is in the time zone where you can't reserve tables (45 minutes before closing) aka 2700 in unix.
func unix_time_in_not_in_closing_time_slot(restaurant_closing_time_to_unix int64, unix_time int64) bool {
	const forty_five_minutes int64 = 2700
	return restaurant_closing_time_to_unix-forty_five_minutes >= unix_time
}

func populate_all_times() []int64 {
	all_times := make([]int64, 0, 96)
	hour := 0
	minutes := 0
	for {
		if hour == 24 {
			break
		}
		all_times = append(all_times, get_unix_from_time(hour, minutes))
		if minutes < 45 || minutes > 0 {
			minutes = minutes + 15
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
func get_time_slots_from_current_point_forward(current_time string) []covered_times {
	all_possible_time_slots := [...]covered_times{
		{time: "0200", time_window_start: "0000", time_window_end: "0600"},
		{time: "0800", time_window_start: "0600", time_window_end: "1200"},
		{time: "1400", time_window_start: "1200", time_window_end: "1800"},
		{time: "2000", time_window_start: "1800", time_window_end: "0000"},
	}
	for time_slot_index, time_slot := range all_possible_time_slots {
		if current_time < time_slot.time_window_end {
			return all_possible_time_slots[time_slot_index:]
		}
	}
	return nil
}

/*
Used to get all the time slots in between the graph start and graph end.
E.g. if start is 2348 and end is 0100, it will get time slots 0000, 0015, 0030, 0045, 0100.
*/
// this throws for some reason.
func time_slots_in_between(start_time string, end_time string, reservation_times []int64) ([]string, error) {
	start_time = convert_uneven_minutes_to_even(start_time)
	end_time = convert_uneven_minutes_to_even(end_time)

	if start_time == "" || end_time == "" {
		return nil, errors.New("error converting uneven minutes to even minutes")
	}
	if start_time == end_time {
		return nil, errors.New("trying to get a time_slot with 2 identical timestamps")
	}
	start_time_minutes, err := strconv.Atoi(start_time[len(start_time)-2:])
	if err != nil {
		return nil, errors.New("error converting start time minutes to int")
	}
	start_time_hours, err := strconv.Atoi(start_time[:len(start_time)-2])
	if err != nil {
		return nil, errors.New("error converting start time hours to int")
	}
	end_time_minutes, err := strconv.Atoi(end_time[len(end_time)-2:])
	if err != nil {
		return nil, errors.New("error converting end time minutes to int")
	}
	end_time_hours, err := strconv.Atoi(end_time[:len(end_time)-2])
	if err != nil {
		return nil, errors.New("error converting end time hours to int")
	}

	start_time_unix := get_unix_from_time(start_time_hours, start_time_minutes)
	end_time_unix := get_unix_from_time(end_time_hours, end_time_minutes)

	if start_time_unix > end_time_unix {
		return nil, errors.New("start_time_unix was larger than end_time_unix")
	}

	var times_we_want []string
	for _, reservation_time := range reservation_times {
		if reservation_time > start_time_unix && reservation_time <= end_time_unix {
			times_we_want = append(times_we_want, get_string_time_from_unix(reservation_time))
		}
	}

	return times_we_want, nil
}
