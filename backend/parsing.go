package main

import (
	"errors"
	"strings"
)

func get_all_possible_reservation_times() *[96]string {
	return &[...]string{
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
func binary_search(a *[96]string, x *string) int {
	r := -1 // not found
	start := 0
	end := len(a) - 1
	for start <= end {
		mid := (start + end) / 2
		if a[mid] == *x {
			r = mid // found
			break
		} else if a[mid] < *x {
			start = mid + 1
		} else if a[mid] > *x {
			end = mid - 1
		}
	}
	return r
}

// this will return all the times in between a certain start time and end time.
func return_time_slots_in_between(start *string, end *string) (*[]string, error) {
	all_possible_reservation_times := get_all_possible_reservation_times()
	start_pos := binary_search(all_possible_reservation_times, start)
	end_pos := binary_search(all_possible_reservation_times, end)
	if start_pos == -1 || end_pos == -1 {
		return nil, errors.New("could not find the corresponding indices from time slot array")
	}
	times_in_between := all_possible_reservation_times[start_pos:end_pos]
	return &times_in_between, nil
}

// Gets the restaurants from the passed in argument. Returns error if nothing is found.
func filter_restaurants_from_city(city *string) (*[]response_fields, error) {
	restaurants := getAllRestaurantsFromRaflaamoApi()
	captured_restaurants := make([]response_fields, 0, len(*restaurants))

	for _, restaurant := range *restaurants {
		if strings.Contains(strings.ToLower(*restaurant.Address.Municipality.Fi_FI), strings.ToLower(*city)) {
			captured_restaurants = append(captured_restaurants, restaurant)
		}
	}
	if len(captured_restaurants) == 0 {
		return &captured_restaurants, errors.New("no restaurants found")
	}
	return &captured_restaurants, nil
}
