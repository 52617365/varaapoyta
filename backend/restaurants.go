package main

import (
	"errors"
	"regexp"
	"strings"
)

// TODO: use goroutines for requests
// TODO: this is too slow when we're doing multiple restaurants
func get_available_tables(restaurants []response_fields, amount_of_eaters int) []response_fields {
	// Getting current_time, so we can avoid checking times from the past.
	current_time := get_current_date_and_time()
	// All possible time slots we need to check, it does not contain time slots from the past.
	time_slots_to_check_from_graph_api := get_graph_time_slots_from_current_point_forward(current_time.time)
	// 11:00, 11:15, 11:30 and so on.
	all_time_intervals := get_all_raflaamo_time_intervals()

	for _, restaurant := range restaurants {
		// If we can't find the id from url, just continue on to the next one because without the id we can't find the reservation page.
		id_from_reservation_page_url, err := get_id_from_reservation_page_url(restaurant)
		if err != nil {
			continue
		}

		restaurant_office_hours := get_opening_and_closing_time_from(restaurant)
		// TODO: pass all_time_intervals into function instead of time_intervals_in_between_office_hours and then do the branch checking inside the function instead of in the previous function.
		available_intervals_from_graph_api, err := get_available_time_intervals_from_graph_api(restaurant_office_hours.opening, restaurant_office_hours.closing, id_from_reservation_page_url, time_slots_to_check_from_graph_api, amount_of_eaters, all_time_intervals, current_time)
		if err != nil {
			continue
		}
		// Here we populate the empty field time slot with all the available time slots.
		// This is expected behavior because we planned on populating it later on.
		restaurant.available_time_slots = available_intervals_from_graph_api
	}
	return restaurants
}

// We do this because the id from the "Id" field is not always the same as the id needed in the reservation page.
func get_id_from_reservation_page_url(restaurant response_fields) (string, error) {
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

// Checks to see if reservation_page_url contains the correct url, sometimes the url is something related to renting a table
// Which will result in an invalid regex match when trying to get id from reservation_page_url.
// @Note this is raflaamo's fault, but we have to deal with it.
func reservation_page_url_is_not_valid(reservation_page_url string) bool {
	return !strings.Contains(reservation_page_url, "https://s-varaukset.fi/online/reservation/fi")
}

// We determine if there is a time slot with open tables by looking at the "color" field in the response.
// The color field will contain "transparent" if it does not contain a graph (open times), else it contains nil (meaning there are open tables)
func time_slot_does_not_contain_open_tables(data parsed_graph_data) bool {
	return data.Intervals[0].Color == "transparent"
}

type restaurant_time struct {
	opening int64
	closing int64
}

func get_opening_and_closing_time_from(restaurant response_fields) restaurant_time {
	// Converting restaurant_start_time to unix, so we can compare it easily.
	restaurant_start_time := get_unix_from_time(restaurant.Openingtime.Restauranttime.Ranges[0].Start)
	// We minus 1 hour from the end time because restaurants don't take reservations before that time slot.
	// IMPORTANT: E.g. if restaurant closes at 22:00, the last possible reservation time is 21:00.
	const one_hour_unix int64 = 3600
	restaurant_ending_time := get_unix_from_time(restaurant.Openingtime.Restauranttime.Ranges[0].End) - one_hour_unix

	return restaurant_time{
		opening: restaurant_start_time,
		closing: restaurant_ending_time,
	}
}
