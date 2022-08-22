package main

import (
	"errors"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
)

// TODO: use goroutines for requests
// TODO: refactor this excuse of a fucking code block, how the fuck did it get here fuck.
func get_available_tables(restaurants []response_fields, current_time date_and_time, amount_of_eaters int) []restaurant_with_available_times_struct {
	re, _ := regexp.Compile(`[^fi/]\d+`) // This regex gets the first number match from the TableReservationLocalized JSON field which is the id we want. https://regex101.com/r/NtFMrz/1
	// All possible time slots we need to check, it does not contain time slots from the past.
	all_possible_time_slots := get_time_slots_from_current_point_forward(current_time.time)

	// This will contain all the available time slots from all restaurants after loop runs.
	all_restaurants_with_available_times := make([]restaurant_with_available_times_struct, 0, 80)

	for _, restaurant := range restaurants {
		id_from_reservation_page_url, err := get_id_from_reservation_page_url(restaurant, re)

		// If we can't find the id from url, just continue on to the next one because then we can't find the reservation page.
		if err != nil {
			continue
		}

		// Here the available_time_slots will be populated once the next for loop iterates all the time_slots.
		single_restaurant_with_available_times := restaurant_with_available_times_struct{
			restaurant:           restaurant,
			available_time_slots: make([]string, 0, len(all_possible_time_slots)),
		}

		// If there is no time ranges available for the restaurant, we just assume it does not even exist.
		if restaurant.Openingtime.Restauranttime.Ranges == nil {
			continue
		}

		// Converting restaurant_start_time to unix, so we can compare it easily.
		restaurant_start_time := get_unix_from_time(restaurant.Openingtime.Restauranttime.Ranges[0].Start)

		// We minus 1 hour from the end time because restaurants don't take reservations before that time slot.
		// IMPORTANT: E.g. if restaurant closes at 22:00, the last possible reservation time is 21:00.
		const one_hour_unix int64 = 3600
		restaurant_ending_time := get_unix_from_time(restaurant.Openingtime.Restauranttime.Ranges[0].End) - one_hour_unix

		all_reservation_times, err := get_all_reservation_times(restaurant_start_time, restaurant_ending_time)
		if err != nil {
			continue
		}

		// Iterating over all possible time slots (0200, 0800, 1400, 2000) to cover the whole 24h window (each time slot covers a 6h window.)
		// However, all all_possible_time_slots does not contain time slots from the past.
		for _, time_slot := range all_possible_time_slots {
			// TODO: store result in channel.
			results_from_graph_api := make(chan string)
			time_slots_from_graph_api, err := get_time_slots_from_graph_api(id_from_reservation_page_url, time_slot.time, amount_of_eaters)
			if err != nil {
				// it's err if there was an error connecting to raflaamo API or if there were no results.
				continue
			}
			// Adding 10800000(ms) to the time to match utc +2 or +3 (finnish time) (10800000 ms corresponds to 3h)
			time_slots_from_graph_api.Intervals[0].From += 10800000
			time_slots_from_graph_api.Intervals[0].To += 10800000

			graph_end_unix := time_slots_from_graph_api.Intervals[0].To

			// time_slows_in_between stores result in channel and contains a non nil value if it had an error.
			is_err := time_slots_in_between(current_time.time, graph_end_unix, results_from_graph_api, all_reservation_times)
			if is_err != nil {
				continue
			}

			// If slice containing time slots does not already contain the time slot, add the time slot from the channel.
			if !slices.Contains(single_restaurant_with_available_times.available_time_slots, <-results_from_graph_api) {
				single_restaurant_with_available_times.available_time_slots = append(single_restaurant_with_available_times.available_time_slots, <-results_from_graph_api)
			}

			// If slice containing time slots does not already contain the time slot, add the time slot from the channel.
			if !slices.Contains(single_restaurant_with_available_times.available_time_slots, <-results_from_graph_api) {
				single_restaurant_with_available_times.available_time_slots = append(single_restaurant_with_available_times.available_time_slots, <-results_from_graph_api)
			}
		}
		all_restaurants_with_available_times = append(all_restaurants_with_available_times, single_restaurant_with_available_times)
	}
	return all_restaurants_with_available_times
}

// We do this because the id from the "Id" field is not always the same as the id needed in the reservation page.
func get_id_from_reservation_page_url(restaurant response_fields, re *regexp.Regexp) (string, error) {
	reservation_page_url := restaurant.Links.TableReservationLocalized.Fi_FI
	if restaurant_does_not_contain_reservation_page(restaurant) {
		return "", errors.New("restaurant did not contain reservation page url")
	}
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

// Some restaurants don't even contain a reservation page url, these restaurants are useless to us, so we make sure to check it.
func restaurant_does_not_contain_reservation_page(restaurant response_fields) bool {
	return len(restaurant.Links.TableReservationLocalized.Fi_FI) == 0
}
