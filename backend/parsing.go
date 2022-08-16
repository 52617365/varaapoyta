package main

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

// FIX: there is a problem where accessing an element results in out of index, this is caused by numbers that when converted go to the start of the array, E.g. 23:49 converts to 0000, resulting in
// the start_time index being larger than the end_time index.
// @Solution: store numbers in the backwards order? (Highest to lowest).
func get_all_possible_reservation_times() [96]string {
	return [...]string{
		"2345", "2330", "2315", "2300",
		"2245", "2230", "2215", "2200",
		"2145", "2130", "2115", "2100",
		"2045", "2030", "2015", "2000",
		"1945", "1930", "1915", "1900",
		"1845", "1830", "1815", "1800",
		"1745", "1730", "1715", "1700",
		"1645", "1630", "1615", "1600",
		"1545", "1530", "1515", "1500",
		"1445", "1430", "1415", "1400",
		"1345", "1330", "1315", "1300",
		"1245", "1230", "1215", "1200",
		"1145", "1130", "1115", "1100",
		"1045", "1030", "1015", "1000",
		"0945", "0930", "0915", "0900",
		"0845", "0830", "0815", "0800",
		"0745", "0730", "0715", "0700",
		"0645", "0630", "0615", "0600",
		"0545", "0530", "0515", "0500",
		"0445", "0430", "0415", "0400",
		"0345", "0330", "0315", "0300",
		"0245", "0230", "0215", "0200",
		"0145", "0130", "0115", "0100",
		"0045", "0030", "0015", "0000",
	}
}

// Binary search algorithm that returns the index of an element in array or -1 if none found.

// In all cases it should find something if we have done the conversion correctly before function call.
func reverse_binary_search(a [96]string, x string) int {
	r := -1 // not found
	//start := 0
	//end := len(a) - 1
	end := 0
	//start := len(a) - 1
	start := len(a) - 1
	//for start <= end {
	for start >= end { // 95 >= 0
		mid := (start + end) / 2 // 95 + 0 / 2
		if a[mid] == x {         // checks if middle is equal to
			r = mid // found
			break
		} else if a[mid] < x { // checks if middle is smaller than the thing we're trying to find
			//start = mid + 1 // start = 47 + 1
			start = mid - 1 // start = 47 - 1
		} else if a[mid] > x { // checks if middle is larger than the thing we're trying to find
			//end = mid - 1 // start = 47 - 1
			end = mid + 1 // start = 47 + 1
		}
	}
	return r
}

// Checks if array contains an element.
func contains(s [5]string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// TODO: if end time gets passed in here, it should not convert to the even number after it since the restaurant could already be closed then.
// returns an even number that is supported by the raflaamo site.
func convert_uneven_minutes_to_even(our_number string) string {
	// Contains all the possible even time slots.
	// 100 is equivalent to E.g. 17:00 (00).
	even_time_slot_minutes := [...]string{"15", "30", "45", "60", "00"}

	our_number_length := len(our_number)

	// Numbers are the last 2 characters in all cases.
	our_number_minutes := our_number[our_number_length-2 : our_number_length]

	our_number_hours := our_number[:our_number_length-2]

	if time_is_already_even(even_time_slot_minutes, our_number_minutes) {
		return our_number
	}
	if our_number_minutes < even_time_slot_minutes[0] {
		even_number := our_number_hours + even_time_slot_minutes[0]
		return even_number
	}
	if our_number_minutes < even_time_slot_minutes[1] {
		even_number := our_number_hours + even_time_slot_minutes[1]
		return even_number
	}
	if our_number_minutes < even_time_slot_minutes[2] {
		even_number := our_number_hours + even_time_slot_minutes[2]
		return even_number
	}
	if our_number_minutes < even_time_slot_minutes[3] {
		// Checking if its 23 to avoid incrementing it to 24 which would be invalid since 24 is represented as 00 (00:00).

		if our_number_hours == "23" {
			return "0000"
		}
		our_number_hour_as_integer, err := strconv.Atoi(our_number_hours)
		if err != nil {
			return ""
		}
		our_number_hour_as_integer++

		// Converting hours back to strings, so we match the original format.

		// Checking if we need to add a 0 before the number because after conversion numbers under 10 will be E.g. "900" when we want it to be "0900".
		// numbers above 10 will not have this problem cuz they will be E.g. "1000" without doing anything.
		if our_number_hour_as_integer < 10 {
			even_number := "0" + strconv.Itoa(our_number_hour_as_integer) + "00"
			return even_number
		}

		even_number := strconv.Itoa(our_number_hour_as_integer) + "00"
		return even_number
	}
	return ""
}

// Checks to see if the time passed in has minutes that we consider even (00, 15, 30, 45).
func time_is_already_even(even_time_slot_minutes [5]string, our_number_minutes string) bool {
	return contains(even_time_slot_minutes, our_number_minutes)
}

// Used to get all the time slots in between the graph start and graph end.
// E.g. if start is 2348 and end is 0100, it will get time slots 0100, 0045, 0030, 0015, 0000.

func return_time_slots_in_between(start string, end string) ([]string, error) {
	all_possible_reservation_times := get_all_possible_reservation_times()
	start_to_even := convert_uneven_minutes_to_even(start)
	end_to_even := convert_uneven_minutes_to_even(end)
	if start_to_even == "" || end_to_even == "" {
		return nil, errors.New("error converting uneven minutes to even minutes")
	}

	start_pos := reverse_binary_search(all_possible_reservation_times, start_to_even)
	end_pos := reverse_binary_search(all_possible_reservation_times, end_to_even)
	if start_pos == -1 || end_pos == -1 {
		return nil, errors.New("could not find the corresponding indices from time slot array")
	}

	// TODO: fix this.
	if end_pos > start_pos {
		log.Fatalln("this should not happen")
	}

	// Here we're checking start the start_pos is not larger than the end_pos because making a slice from that range is going to result in a panic.
	// (It's something related to the sequence of times in the return value of get_all_possible_reservation_times)
	// @Solution, we're going to store the times in the reverse order.

	times_in_between := all_possible_reservation_times[end_pos:start_pos]
	return times_in_between, nil
}

// Gets the restaurants from the passed in argument. Returns error if nothing is found.
func filter_restaurants_from_city(city string) ([]response_fields, error) {
	restaurants := getAllRestaurantsFromRaflaamoApi()
	captured_restaurants := make([]response_fields, 0, len(*restaurants))

	for _, restaurant := range *restaurants {
		if strings.Contains(strings.ToLower(restaurant.Address.Municipality.Fi_FI), strings.ToLower(city)) {
			captured_restaurants = append(captured_restaurants, restaurant)
		}
	}
	if len(captured_restaurants) == 0 {
		return captured_restaurants, errors.New("no restaurants found")
	}
	return captured_restaurants, nil
}
