package main

import (
	"errors"
	"strings"

	"golang.org/x/exp/slices"
)

// TODO: Even after our recent changes, this function is unfortunately still too slow (10 seconds for helsinki restaurants), we NEED to optimize it. Now is the time.
// TODO: this currently returns an empty body.
func get_available_tables(city string, amount_of_eaters int) ([]response_fields, error) {
	// Getting current_time, so we can avoid checking times from the past.
	current_time := get_current_date_and_time()
	all_time_intervals := get_all_raflaamo_time_intervals()
	time_slots_to_check := get_graph_time_slots_from_current_point_forward(current_time.time)

	raflaamo, err := init_restaurants_api()
	if err != nil {
		return []response_fields{}, errors.New("failed to construct the http client, contact the developer")
	}
	response, err := raflaamo.get()
	if err != nil {
		return []response_fields{}, errors.New("raflaamo api most likely down")
	}

	jobs := make(chan job)
	restaurants_with_opening_times := make([]response_fields, 50)

	for _, restaurant := range response {
		if restaurant_format_is_incorrect(city, restaurant) {
			continue
		}
		restaurant_additional_information, err := init_additional_information(restaurant)
		if err != nil {
			continue
		}
		restaurant_additional_information.add()

		restaurant_id := restaurant_additional_information.restaurant.Links.TableReservationLocalizedId
		if restaurant_id == "" {
			// restaurant id could not be found.
			continue
		}

		kitchen_times, err := get_opening_and_closing_time_from_kitchen_time(restaurant)
		if err != nil {
			continue
		}

		for range time_slots_to_check {
			go worker(jobs, restaurant.response)
		}
		spawn_jobs(jobs, restaurant_id, kitchen_times, time_slots_to_check, restaurant, amount_of_eaters, current_time, city)

		for graph_response := range restaurant.response {
			available_time_slots, err := extract_available_time_intervals_from_response(graph_response, current_time, kitchen_times, all_time_intervals)
			if err != nil {
				continue
			}
			// Avoiding duplicates.
			for _, available_time_slot := range available_time_slots {
				if slices.Contains(restaurant.Available_time_slots, available_time_slot) {
					restaurant.Available_time_slots = append(restaurant.Available_time_slots, available_time_slot)
				}
			}
			restaurants_with_opening_times = append(restaurants_with_opening_times, restaurant)
		}
	}
	return restaurants_with_opening_times, nil
}

// Checks to see if reservation_page_url contains the correct url, sometimes the url is something related to renting a table
// Which will result in an invalid regex match when trying to get id from reservation_page_url.
// @Note this is raflaamo's fault, but we have to deal with it.
func reservation_page_url_is_not_valid(reservation_page_url string) bool {
	return !strings.Contains(reservation_page_url, "https://s-varaukset.fi/online/reservation/fi")
}

type restaurant_time struct {
	opening int64
	closing int64
}

func restaurant_format_is_incorrect(city string, restaurant response_fields) bool {
	// Converting to lower so that we don't run into problems when comparing it.
	city = strings.ToLower(city)

	// The restaurant times are nil if there are no time ranges available, we want to skip those because they are useless to us without times.
	if restaurant.Openingtime.Restauranttime.Ranges == nil || restaurant.Openingtime.Kitchentime.Ranges == nil {
		return true
	}
	if restaurant.Openingtime.Kitchentime.Ranges[0].Start == "" || restaurant.Openingtime.Kitchentime.Ranges[0].End == "" {
		return true
	}
	// It sometimes has this weird (MISSING) thing and that is an indication that there is something scary happening with the reservation page link so we avoid those.
	if strings.Contains(restaurant.Links.TableReservationLocalized.Fi_FI, "(MISSING)") {
		return true
	}
	// We check that the city is the correct one since we want to filter them out depending on the passed in city.
	if strings.ToLower(restaurant.Address.Municipality.Fi_FI) != city {
		return true
	}
	// Reservation url is sometimes empty, skip at this point cuz we can't reserve a table without a reservation page url.
	if restaurant.Links.TableReservationLocalized.Fi_FI == "" {
		return true
	}
	restaurant_office_hours, err := get_opening_and_closing_time_from_kitchen_time(restaurant)
	if err != nil {
		return true
	}
	// Checking to see if the timestamps are fucked here, so we don't have to check them later.
	// We have already checked that the ranges exist in the previous condition (restaurant.Openingtime.Restauranttime.Ranges != nil)
	is_invalid_date := restaurant_office_hours.opening >= restaurant_office_hours.closing
	return is_invalid_date
}
