package main

import (
	"errors"
	"regexp"
)

type additional_information struct {
	restaurant    response_fields
	time_slots    chan Result
	kitchen_times restaurant_time
}

func init_additional_information(restaurant response_fields, time_slots_to_check_length int) additional_information {
	kitchen_office_hours := get_opening_and_closing_time_from_kitchen_time(restaurant)
	return additional_information{
		restaurant:    restaurant,
		kitchen_times: kitchen_office_hours,
		time_slots:    make(chan Result, time_slots_to_check_length),
	}
}

/*
This struct exists because it lets us filter out the restaurants we're not interested in (E.g. from city we didn't want)
whilst associating the response with the correct restaurant.

also,

Stores certain additional information about the restaurant that we might be interested in.
For example:
- The id from the reservation page url
- Time till restaurant closes.
- Time till restaurants kitchen closes.
*/
func (add *additional_information) add() error {
	restaurant_id, err := add.get_id_from_reservation_page_url()
	if err != nil {
		return err
	}
	add.restaurant.Links.TableReservationLocalizedId = restaurant_id
	time := time_utils{
		current_time: get_current_date_and_time(),
		closing_time: add.kitchen_times.closing,
	}

	time_till_restaurant_closed := time.get_time_till_restaurant_closing_time()
	time_till_restaurants_kitchen_closed := time.get_time_left_to_reserve()

	add.restaurant.Openingtime.Time_till_restaurant_closed_hours = time_till_restaurant_closed.hour
	add.restaurant.Openingtime.Time_till_restaurant_closed_minutes = time_till_restaurant_closed.minutes

	add.restaurant.Openingtime.Time_left_to_reserve_hours = time_till_restaurants_kitchen_closed.hour
	add.restaurant.Openingtime.Time_left_to_reserve_minutes = time_till_restaurants_kitchen_closed.minutes
	return nil
}

/*
Gets the id from a restaurants reservation url.
We have already checked here that the string we're matching contains the string "https://s-varaukset.fi/online/reservation/fi", so it should return something.
We return an error in case but in reality it really should not return error.
*/
func (add additional_information) get_id_from_reservation_page_url() (string, error) {
	restaurant := add.restaurant
	re, _ := regexp.Compile(`[^fi/]\d+`) // This regex gets the first number match from the TableReservationLocalized JSON field which is the id we want. https://regex101.com/r/NtFMrz/1
	reservation_page_url := restaurant.Links.TableReservationLocalized.Fi_FI

	id_from_reservation_page_url := re.FindString(reservation_page_url)

	// If regex could not match or if url was invalid (happens sometimes cuz API is weird).
	if id_from_reservation_page_url == "" {
		return "", errors.New("regex did not match anything, something wrong with reservation_page_url")
	}
	return id_from_reservation_page_url, nil
}
