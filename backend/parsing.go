package main

import (
	"errors"
	"strconv"
	"strings"
)

func is_not_valid_format(our_number string) bool {
	if _, err := strconv.ParseInt(our_number, 10, 64); err != nil {
		return true
	}
	if len(our_number) != 4 {
		return true
	}
	if our_number == "" {
		return true
	}
	return false
}

// Gets the restaurants from the passed in argument. Returns error if nothing is found.
func filter_restaurants_from_city(city string) ([]response_fields, error) {
	restaurants, err := get_all_restaurants_from_raflaamo_api()
	if err != nil {
		return nil, errors.New("there was an error connecting to the raflaamo api")
	}
	captured_restaurants := make([]response_fields, 0, len(restaurants))

	for _, restaurant := range restaurants {
		if strings.Contains(strings.ToLower(restaurant.Address.Municipality.Fi_FI), strings.ToLower(city)) {
			captured_restaurants = append(captured_restaurants, restaurant)
		}
	}
	return captured_restaurants, nil
}
