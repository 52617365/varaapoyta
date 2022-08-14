package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// This file handles everything related to parsing shit.

// Numbers after 1000 are 4 digits so check if number is under 1000, if so, add trailing zero.
func convert_times_to_string(times []int) []string {
	formattedStrings := make([]string, len(times))
	// This for loop first makes sure that everything is the same length (converts 1000< to 4 digits E.g. 800 is 0800)
	// This is done to handle them all the same way.
	for i, t := range times {
		var formattedString string
		// AKA if number length is only 3
		if t < 1000 {
			formattedString = fmt.Sprintf("0%d", t)
			// Here formattedString will look like "1230" so we can assume it will always have 3 indices.
		} else {
			// AKA if number length is already 4
			formattedString = strconv.Itoa(t)
		}
		hour := formattedString[:2]
		minutes := formattedString[2:]
		formattedString = fmt.Sprintf("%s:%s", hour, minutes)
		formattedStrings[i] = formattedString
	}
	return formattedStrings
}

func getAllPossibleTimes() *[]string {
	// Here we have all the possible all_possible_reservation_times when you can reserve a table.
	all_possible_reservation_times := get_all_possible_reservation_times()

	// Won't be an error since getCurrentTime returns right value everytime.
	currentTime, _ := strconv.Atoi(strings.ReplaceAll(getCurrentTime(), ":", ""))
	var timesWeWant []string

	// Get all the times we want. (List is sorted, so we can assume that if a number is larger, everything after it will be too, so we don't need a branch for everything after that)
	for i := 0; i < len(all_possible_reservation_times); i++ {
		if all_possible_reservation_times[i] > currentTime {
			timesWeWant = convert_times_to_string(all_possible_reservation_times[i:])
			break
		}
	}
	return &timesWeWant
}

func get_all_possible_reservation_times() [63]int {
	all_possible_reservation_times := [...]int{
		800, 815, 830, 845, 900, 915, 930, 945, 1000, 1015, 1030,
		1100, 1115, 1130, 1145, 1200, 1215, 1230, 1245, 1300,
		1315, 1330, 1345, 1400, 1415, 1430, 1445, 1500, 1515, 1530,
		1545, 1600, 1615, 1630, 1645, 1700, 1715, 1730, 1745, 1800,
		1815, 1830, 1845, 1900, 1915, 1930, 1945, 2000, 2015, 2030,
		2045, 2100, 2115, 2130, 2145, 2200, 2215, 2230, 2245, 2300,
		2315, 2330, 2345,
	}
	return all_possible_reservation_times
}

// Binary search algorithm that returns the index of an element in array or -1 if none found.
func binary_search(a [63]int, x int) int {
	r := -1 // not found
	start := 0
	end := len(a) - 1
	for start <= end {
		mid := (start + end) / 2
		if a[mid] == x {
			r = mid // found
			break
		} else if a[mid] < x {
			start = mid + 1
		} else if a[mid] > x {
			end = mid - 1
		}
	}
	return r
}

// this will return all the times in between a certain start time and end time.
func return_time_slots_in_between(start int, end int) (*[]int, error) {
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
