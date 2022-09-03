package main

import (
	"errors"
	"strconv"
	"strings"
)

type relative_time struct {
	hour    int
	minutes int
}

func get_time_till_restaurant_closing_time(closing_time int64) relative_time {
	// TODO: figure out why closing time is 30 minutes too low, it's because of we take kitchen time now instead of restaurant time.
	// we minused one hour from it cuz they don't take reservations in that time slot, but they're still technically open, so we add it back here, this is the only place where we add it back.
	const one_hour_unix int64 = 3600
	closing_time += one_hour_unix
	current_time := get_current_date_and_time()
	// already closed.
	if closing_time <= current_time.time {
		return relative_time{hour: -1, minutes: -1}
	}

	time_left_to_closing_unix := closing_time - current_time.time

	relative_time_string := get_string_time_from_unix(time_left_to_closing_unix)
	relative_time_string = strings.Replace(relative_time_string, ":", "", -1)

	if is_not_valid_format(relative_time_string) {
		return relative_time{hour: -1, minutes: -1}
	}

	minutes, _ := strconv.Atoi(relative_time_string[len(relative_time_string)-2:])
	hour, _ := strconv.Atoi(relative_time_string[:len(relative_time_string)-2])

	return relative_time{
		hour:    hour,
		minutes: minutes,
	}

}
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
			restaurant_office_hours := get_opening_and_closing_time_from_kitchen_time(restaurant)
			// Checking to see if the timestamps are fucked here, so we don't have to check them later.
			// We have already checked that the ranges exist in the previous condition (restaurant.Openingtime.Restauranttime.Ranges != nil)
			if !(restaurant_office_hours.opening >= restaurant_office_hours.closing) {
				captured_restaurants = append(captured_restaurants, restaurant)
			}
		}
	}
	return captured_restaurants, nil
}
