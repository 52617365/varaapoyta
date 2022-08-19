package main

import (
	"errors"
	"regexp"
	"strings"
)

// TODO: use goroutines for requests
func getAvailableTables(restaurants []response_fields, amount_of_eaters int) []restaurant_with_available_times_struct {
	re, _ := regexp.Compile(`[^fi/]\d+`) // This regex gets the first number match from the TableReservationLocalized JSON field which is the id we want. https://regex101.com/r/NtFMrz/1

	current_date := get_current_date_and_time()

	all_possible_time_slots := get_time_slots_from_current_point_forward(current_date.time)

	// There can be maximum of restaurants * all_possible_time_slots, so we allocate the worst case scenario here to avoid reallocation's.
	total_memory_to_reserve_for_all_restaurant_time_slots := len(restaurants) * len(all_possible_time_slots)

	// This will contain all the available time slots from all restaurants after loop runs.
	all_restaurants_with_available_times := make([]restaurant_with_available_times_struct, 0, total_memory_to_reserve_for_all_restaurant_time_slots)

	for _, restaurant := range restaurants {
		id_from_reservation_page_url, err := get_id_from_reservation_page_url(restaurant, re)

		if err != nil {
			continue
		}

		// Here the available_time_slots will be populated once the next for loop iterates all the time_slots.
		single_restaurant_with_available_times := restaurant_with_available_times_struct{
			restaurant:           restaurant,
			available_time_slots: make([]string, 0, len(all_possible_time_slots)),
		}

		// Iterating over all possible time slots (0200, 0800, 1400, 2000) to cover the whole 24h window (each time slot covers a 6h window.)
		for _, time_slot := range all_possible_time_slots {
			time_slots_from_graph_api, err := get_time_slots_from_graph_api(id_from_reservation_page_url, current_date.date, time_slot.time, amount_of_eaters)
			if err != nil {
				continue
			}

			// At this point in the code we have already made all the necessary checks to confirm that a graph is visible for a time slot, and we can extract information from it.

			unix_timestamp_struct_of_available_table := convert_unix_timestamp_to_finland_time(time_slots_from_graph_api)

			restaurant_starting_time := restaurant.Openingtime.Restauranttime.Ranges[0].Start
			restaurant_closing_time := restaurant.Openingtime.Restauranttime.Ranges[0].End

			all_reservation_times, err := get_all_reservation_times(restaurant_starting_time, restaurant_closing_time) // in reality, it's not all because we need to consider restaurants closing time.
			if err != nil {
				continue
			}
			time_slots, err := time_slots_in_between(current_date.time, unix_timestamp_struct_of_available_table.time_window_end, all_reservation_times)

			if err != nil {
				continue
			}

			single_restaurant_with_available_times.available_time_slots = append(single_restaurant_with_available_times.available_time_slots, time_slots...)
		}
		// Here after iterating over all time slots for the restaurant, we store the results.
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
func time_slot_does_not_contain_open_tables(data *parsed_graph_data) bool {
	return data.Intervals[0].Color == "transparent"
}

// Some restaurants don't even contain a reservation page url, these restaurants are useless to us so we make sure to check it.
func restaurant_does_not_contain_reservation_page(restaurant response_fields) bool {
	return len(restaurant.Links.TableReservationLocalized.Fi_FI) == 0
}
