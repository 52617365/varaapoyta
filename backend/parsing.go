package main

import (
	"errors"
	"strconv"
	"strings"
)

// Returns an even number that is supported by the raflaamo site.
func convert_uneven_minutes_to_even(our_number string) string {
	if len(our_number) < 4 {
		return ""
	}
	our_number_length := len(our_number)

	// Numbers are the last 2 characters in all cases.
	our_number_minutes := our_number[our_number_length-2 : our_number_length]

	our_number_hours := our_number[:our_number_length-2]

	if time_is_already_even(our_number_minutes) {
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

		/*
			Checking if we need to add a 0 before the number because after conversion numbers under 10 will be
			E.g. "900" when we want it to be "0900".
			Numbers above 10 will not have this problem cuz they will be E.g. "1000" without doing anything.
		*/

		// Converting hours back to strings, so we match the original format.
		if our_number_hour_as_integer < 10 {
			our_number_hour_as_integer++
			even_number := "0" + strconv.Itoa(our_number_hour_as_integer) + "00"
			return even_number
		}

		// Converting hours back to strings, so we match the original format.
		even_number := strconv.Itoa(our_number_hour_as_integer) + "00"
		return even_number
	}
	return ""
}

// Checks to see if the time passed in has minutes that we consider even (00, 15, 30, 45).
func time_is_already_even(our_number_minutes string) bool {
	if our_number_minutes == "45" || our_number_minutes == "30" || our_number_minutes == "15" || our_number_minutes == "00" {
		return true
	}
	return false
}

// Gets the restaurants from the passed in argument. Returns error if nothing is found.
func filter_restaurants_from_city(city string) ([]response_fields, error) {
	restaurants := get_all_restaurants_from_raflaamo_api()
	captured_restaurants := make([]response_fields, 0, len(restaurants))

	for _, restaurant := range restaurants {
		if strings.Contains(strings.ToLower(restaurant.Address.Municipality.Fi_FI), strings.ToLower(city)) {
			captured_restaurants = append(captured_restaurants, restaurant)
		}
	}
	if len(captured_restaurants) == 0 {
		return captured_restaurants, errors.New("no restaurants found")
	}
	return captured_restaurants, nil
}
