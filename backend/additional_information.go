package main

import (
	"errors"
	"regexp"
)

type additional_information struct {
	restaurant           response_fields
	kitchen_office_hours restaurant_time
}

func init_additional_information(restaurant response_fields) (additional_information, error) {
	kitchen_office_hours, err := get_opening_and_closing_time_from_kitchen_time(restaurant)
	if err != nil {
		return additional_information{}, err
	}
	return additional_information{
		restaurant:           restaurant,
		kitchen_office_hours: kitchen_office_hours,
	}, nil
}

// Storing reservation id, so we can easily use it later without parsing the reservation url again.
// Adding relative times, so we can display them as a countdown on the page later.
// Adding the relative time to when the restaurant itself closes (Timestamp can be different, then the kitchen time).
// Adding the relative time to when the restaurants kitchen closes (Timestamp can be different, then the restaurant itself).
func (add *additional_information) add() {
	restaurant_id, err := add.get_id_from_reservation_page_url()
	if err == nil {
		add.restaurant.Links.TableReservationLocalizedId = restaurant_id
	}
	time := time_utils{
		current_time: get_current_date_and_time(),
		closing_time: add.kitchen_office_hours.closing,
	}

	time_till_restaurant_closed := time.get_time_till_restaurant_closing_time()
	time_till_restaurants_kitchen_closed := time.get_time_left_to_reserve()

	add.restaurant.Openingtime.Time_till_restaurant_closed_hours = time_till_restaurant_closed.hour
	add.restaurant.Openingtime.Time_till_restaurant_closed_minutes = time_till_restaurant_closed.minutes

	add.restaurant.Openingtime.Time_left_to_reserve_hours = time_till_restaurants_kitchen_closed.hour
	add.restaurant.Openingtime.Time_left_to_reserve_minutes = time_till_restaurants_kitchen_closed.minutes
}

// Returns the id from a reservation page url associated with a restaurant because the id that comes with the restaurant might not be the same as the one in the
// reservation page url and that id is the one that lets us get access to the information related to reserving stuff etc.
func (add additional_information) get_id_from_reservation_page_url() (string, error) {
	restaurant := add.restaurant
	re, _ := regexp.Compile(`[^fi/]\d+`) // This regex gets the first number match from the TableReservationLocalized JSON field which is the id we want. https://regex101.com/r/NtFMrz/1
	reservation_page_url := restaurant.Links.TableReservationLocalized.Fi_FI

	if reservation_page_url_is_not_valid(reservation_page_url) {
		return "", errors.New("reservation_page_url_is_not_valid")
	}
	id_from_reservation_page_url := re.FindString(reservation_page_url)

	// If regex could not match or if url was invalid (happens sometimes cuz API is weird).
	if id_from_reservation_page_url == "" {
		return "", errors.New("regex did not match anything, something wrong with reservation_page_url")
	}
	return id_from_reservation_page_url, nil
}
