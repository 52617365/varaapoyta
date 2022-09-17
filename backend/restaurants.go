package main

import (
	"errors"
	"strings"

	"github.com/gammazero/workerpool"
	"golang.org/x/exp/slices"
)

// TODO: Even after our recent changes, this function is unfortunately still too slow (10 seconds for helsinki restaurants), we NEED to optimize it. Now is the time.
// TODO: this currently returns an empty body.
func get_available_tables(city string, amount_of_eaters int) []response_fields {
	// Using channel because this function would block for around ~2 seconds otherwise.
	raflaamo_api_response := make(chan []response_fields)
	raflaamo_api_response_error := make(chan error)
	go get_all_restaurants_from_raflaamo_api(raflaamo_api_response, raflaamo_api_response_error)

	// Getting current_time, so we can avoid checking times from the past.
	current_time := get_current_date_and_time()
	all_time_intervals := get_all_raflaamo_time_intervals()
	wp := workerpool.New(2)
	time_slots_to_check := get_graph_time_slots_from_current_point_forward(current_time.time)

	if <-raflaamo_api_response_error != nil {
		return nil
	}
	raflaamo_api_restaurants := <-raflaamo_api_response

	for _, restaurant := range raflaamo_api_restaurants {

		if restaurant_format_is_incorrect(city, restaurant) {
			continue
		}
		kitchen_office_hours, err := get_opening_and_closing_time_from_kitchen_time(restaurant)
		if err != nil {
			continue
		}
		restaurant_additional_information := additional_information{
			restaurant:           restaurant,
			kitchen_office_hours: kitchen_office_hours,
		}

		restaurant_id, err := restaurant_additional_information.get_id_from_reservation_page_url()

		if err != nil {
			continue
		}

		restaurant_additional_information.add()

		for _, slot := range time_slots_to_check {
			wp.Submit(func() {
				graph_data, err := interact_with_api(slot, restaurant_id, current_time, amount_of_eaters)
				restaurant.Api_response.err = err
				restaurant.Api_response.response = graph_data
			})
		}

	}
	wp.StopWait()

	restaurants_with_opening_times := make([]response_fields, 50)
	for _, restaurant := range raflaamo_api_restaurants {
		if restaurant_format_is_incorrect(city, restaurant) {
			continue
		}
		if restaurant.Api_response.err != nil {
			continue
		}
		response := restaurant.Api_response.response
		kitchen_times, err := get_opening_and_closing_time_from_kitchen_time(restaurant)
		if err != nil {
			continue
		}
		available_time_slots, err := extract_available_time_intervals_from_response(response, current_time, kitchen_times, all_time_intervals)
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
	return restaurants_with_opening_times
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

func get_opening_and_closing_time_from_kitchen_time(restaurant response_fields) (restaurant_time, error) {
	// Converting restaurant_kitchen_start_time to unix, so we can compare it easily.
	// restaurant_kitchen_start_time := get_unix_from_time(restaurant.Openingtime.Restauranttime.Ranges[0].Start)
	if len(restaurant.Openingtime.Kitchentime.Ranges) == 0 {
		return restaurant_time{}, errors.New("no ranges found")
	}
	restaurant_kitchen_start_time := get_unix_from_time(restaurant.Openingtime.Kitchentime.Ranges[0].Start)
	// We minus 1 hour from the end time because restaurants don't take reservations before that time slot.
	// IMPORTANT: E.g. if restaurant closes at 22:00, the last possible reservation time is 21:00.
	const one_hour_unix int64 = 3600
	// restaurant_kitchen_ending_time := get_unix_from_time(restaurant.Openingtime.Restauranttime.Ranges[0].End) - one_hour_unix
	restaurant_kitchen_ending_time := get_unix_from_time(restaurant.Openingtime.Kitchentime.Ranges[0].End) - one_hour_unix

	return restaurant_time{
		opening: restaurant_kitchen_start_time,
		closing: restaurant_kitchen_ending_time,
	}, nil
}

func restaurant_format_is_incorrect(city string, restaurant response_fields) bool {
	// Converting to lower so that we don't run into problems when comparing it.
	city = strings.ToLower(city)

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
	// The restaurant times are nil if there are no time ranges available, we want to skip those because they are useless to us without times.
	if restaurant.Openingtime.Restauranttime.Ranges == nil || restaurant.Openingtime.Kitchentime.Ranges == nil {
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
	return restaurant_office_hours.opening >= restaurant_office_hours.closing
}
