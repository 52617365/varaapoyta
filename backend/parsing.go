package main

import (
	"errors"
	"strconv"
	"strings"
)

// FIX: there is a problem where accessing an element results in out of index, this is caused by numbers that when converted go to the start of the array, E.g. 23:49 converts to 0000, resulting in
// the start_time index being larger than the end_time index.
func get_all_possible_reservation_times() [96]string {
	return [...]string{
		"0000", "0015", "0030", "0045",
		"0100", "0115", "0130", "0145",
		"0200", "0215", "0230", "0245",
		"0300", "0315", "0330", "0345",
		"0400", "0415", "0430", "0445",
		"0500", "0515", "0530", "0545",
		"0600", "0615", "0630", "0645",
		"0700", "0715", "0730", "0745",
		"0800", "0815", "0830", "0845",
		"0900", "0915", "0930", "0945",
		"1000", "1015", "1030", "1045",
		"1100", "1115", "1130", "1145",
		"1200", "1215", "1230", "1245",
		"1300", "1315", "1330", "1345",
		"1400", "1415", "1430", "1445",
		"1500", "1515", "1530", "1545",
		"1600", "1615", "1630", "1645",
		"1700", "1715", "1730", "1745",
		"1800", "1815", "1830", "1845",
		"1900", "1915", "1930", "1945",
		"2000", "2015", "2030", "2045",
		"2100", "2115", "2130", "2145",
		"2200", "2215", "2230", "2245",
		"2300", "2315", "2330", "2345",
	}
}

// Binary search algorithm that returns the index of an element in array or -1 if none found.

// In all cases it should find something if we have done the conversion correctly before function call.
func binary_search(a [96]string, x string) int {
	r := -1 // not found
	start := 0
	end := len(a) - 1
	for start <= end {
		mid := (start + end) / 2
		if a[mid] == x { // checks if middle is equal to
			r = mid // found
			break
		} else if a[mid] < x { // checks if middle is smaller than the thing we're trying to find
			start = mid + 1
		} else if a[mid] > x { // checks if middle is larger than the thing we're trying to find
			end = mid - 1
		}
	}
	return r
}

/*
	TODO: Check if the restaurants closing time and avoid checks that are in the last hour of the restaurant closing.
	This is because the restaurants don't take table reservations in the last hour.
*/

// Returns an even number that is supported by the raflaamo site.
func convert_uneven_minutes_to_even(our_number string) string {
	// Contains all the possible even time slots.
	// 100 is equivalent to E.g. 17:00 (00).
	our_number_length := len(our_number)

	// Numbers are the last 2 characters in all cases.
	our_number_minutes := our_number[our_number_length-2 : our_number_length]

	our_number_hours := our_number[:our_number_length-2]

	if time_is_already_even([...]string{"15", "30", "45", "00"}, our_number_minutes) {
		return our_number
	}
	if our_number_minutes < "15" {
		even_number := our_number_hours + "15"
		return even_number
	}
	if our_number_minutes < "30" {
		even_number := our_number_hours + "30"
		return even_number
	}
	if our_number_minutes < "45" {
		even_number := our_number_hours + "45"
		return even_number
	}
	if our_number_minutes > "45" {
		// Checking if its 23 to avoid incrementing it to 24 which would be invalid since 24 is represented as 00 (00:00).

		if our_number_hours == "23" {
			return "0000"
		}

		// Converting to integer so we can increment it.
		our_number_hour_as_integer, err := strconv.Atoi(our_number_hours)
		if err != nil {
			return ""
		}

		// Converting hours back to strings, so we match the original format.

		/*
			Checking if we need to add a 0 before the number because after conversion numbers under 10 will be
			E.g. "900" when we want it to be "0900".
			Numbers above 10 will not have this problem cuz they will be E.g. "1000" without doing anything.
		*/
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
func time_is_already_even(even_time_slot_minutes [4]string, our_number_minutes string) bool {
	for _, even_time_slot := range even_time_slot_minutes {
		if even_time_slot == our_number_minutes {
			return true
		}
	}
	return false
}

/*
Used to get all the time slots in between the graph start and graph end.
E.g. if start is 2348 and end is 0100, it will get time slots 0000, 0015, 0030, 0045, 0100.
*/
func return_time_slots_in_between(start string, end string) ([]string, error) {
	all_possible_reservation_times := get_all_possible_reservation_times()
	start_to_even := convert_uneven_minutes_to_even(start)
	end_to_even := convert_uneven_minutes_to_even(end)
	if start_to_even == "" || end_to_even == "" {
		return nil, errors.New("error converting uneven minutes to even minutes")
	}

	start_pos := binary_search(all_possible_reservation_times, start_to_even)
	end_pos := binary_search(all_possible_reservation_times, end_to_even)
	if start_pos == -1 || end_pos == -1 {
		return nil, errors.New("could not find the corresponding indices from time slot array")
	}

	/*
		Handle the case where end_pos is E.g. 23:49 and therefore converts to 0000, goes back to index 0
		which will result in the index in start_pos being bigger than the index in end_pos
		resulting in out of index error.
	*/
	if end_pos < start_pos {
		times_till_end := all_possible_reservation_times[start_pos:]
		times_from_start := all_possible_reservation_times[:end_pos+1]

		space_to_allocate := len(times_from_start) + len(times_till_end)

		times_in_between := make([]string, 0, space_to_allocate)

		times_in_between = append(times_in_between, times_from_start...)
		times_in_between = append(times_in_between, times_till_end...)

		return times_in_between, nil
	}

	times_in_between := all_possible_reservation_times[start_pos:end_pos]
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
