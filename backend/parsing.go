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
func filter_valid_restaurants_from_city(city string) ([]response_fields, error) {
	city = strings.ToLower(city)
	if city == "" {
		return nil, errors.New("no city provided")
	}
	restaurants, err := get_all_restaurants_from_raflaamo_api()
	if err != nil {
		return nil, errors.New("there was an error connecting to the raflaamo api")
	}
	captured_restaurants := make([]response_fields, 0, 30)

	// TODO: shall we filter all the restaurants out here that don't have a reservation link or ranges?
	for _, restaurant := range restaurants {

		// If there is no time ranges available for the restaurant, we just assume it does not even exist.
		// Also, if there is no reservation link the restaurant is useless to us.
		if strings.ToLower(restaurant.Address.Municipality.Fi_FI) == city && restaurant.Openingtime.Restauranttime.Ranges != nil && restaurant.Links.TableReservationLocalized.Fi_FI != "" {
			captured_restaurants = append(captured_restaurants, restaurant)
		}
	}
	return captured_restaurants, nil
}
