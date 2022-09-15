package main

import (
	"errors"
	"regexp"
	"strings"
)

// TODO: Even after our recent changes, this function is unfortunately still too slow (10 seconds for helsinki restaurants), we NEED to optimize it. Now is the time.
// @Add to above: The worker pool implementation is fucked only use it for requests.
// TODO: this currently returns an empty body.
func get_available_tables(city string, amount_of_eaters int) []response_fields {
	// Using channel because this function would block for around ~2 seconds otherwise.
	raflaamo_api_response := make(chan []response_fields)
	raflaamo_api_response_error := make(chan error)

	go get_all_restaurants_from_raflaamo_api(raflaamo_api_response, raflaamo_api_response_error)

	// Getting current_time, so we can avoid checking times from the past.
	current_time := get_current_date_and_time()

	time_slots_to_check_from_graph_api := get_graph_time_slots_from_current_point_forward(current_time.time)

	all_time_intervals := get_all_raflaamo_time_intervals()

	// We capture the outcomes here to avoid blocking for the duration of the request of get_all_restaurants_from_raflaamo_api.
	if <-raflaamo_api_response_error != nil {
		return nil
	}
	raflaamo_api_restaurants := <-raflaamo_api_response
	restaurants_from_city := make(chan response_fields)
	for _, restaurant := range raflaamo_api_restaurants {
		if restaurant_format_is_incorrect(city, restaurant) {
			continue
		}

		// The reservation page urls are sometimes incorrect, so we just skip over that specific restaurant in that case.
		id_from_reservation_page_url, err := get_id_from_reservation_page_url(restaurant)
		if err != nil {
			continue
		}

		// Initializing these into buffered channels, so they don't block.
		restaurant.Graph_api_response = make(chan parsed_graph_data, len(time_slots_to_check_from_graph_api))
		restaurant.Graph_api_response_err = make(chan error, len(time_slots_to_check_from_graph_api))

		kitchen_office_hours := get_opening_and_closing_time_from_kitchen_time(restaurant)
		// Here we add some fields into the restaurant struct which we already have and will be needing in the future.
		restaurant = add_additional_fields(restaurant, id_from_reservation_page_url, kitchen_office_hours)

		// Looping over each time_slot we need to check, in each restaurant and spawning a job for each one.
		// The results will be stored in the restaurant struct itself, this is because we want to associate the response with the restaurant later.
		for _, time_slot := range time_slots_to_check_from_graph_api {
			go func(time_slot covered_times) {
				interact_with_api(restaurants_from_city, restaurant, time_slot, id_from_reservation_page_url, current_time, amount_of_eaters)
			}(time_slot)
		}
		close(restaurant.Graph_api_response)
		close(restaurant.Graph_api_response_err)
	}
	close(restaurants_from_city)
	// wp.StopWait()

	restaurants_with_opening_times := make([]response_fields, 40)
	for restaurant := range restaurants_from_city {
		// The restaurant contains two buffered channels with a length of len(time_slots_to_check_from_graph_api) so we loop over them all.
		for range time_slots_to_check_from_graph_api {
			is_err := <-restaurant.Graph_api_response_err
			if is_err != nil {
				continue
			}
			api_response := <-restaurant.Graph_api_response
			kitchen_office_hours := get_opening_and_closing_time_from_kitchen_time(restaurant)
			available_time_slots, err := extract_available_time_intervals_from_response(api_response, current_time, kitchen_office_hours, all_time_intervals)
			// If this is err, the restaurant is already closed, and therefore we don't want it.
			if err != nil {
				continue
			}
			// Capturing all the available time slots into a shell variable which we made for this purpose.
			// We add it here instead of in add_additional_fields because we only just now know the times, not when we were assigning the other additional fields.
			restaurant.Available_time_slots = append(restaurant.Available_time_slots, available_time_slots...)

			// Saving the restaurant for later.
			restaurants_with_opening_times = append(restaurants_with_opening_times, restaurant)
		}
	}
	return restaurants_with_opening_times
}

// Storing reservation id, so we can easily use it later without parsing the reservation url again.
// Adding relative times, so we can display them as a countdown on the page later.
// Adding the relative time to when the restaurant itself closes (Timestamp can be different, then the kitchen time).
// Adding the relative time to when the restaurants kitchen closes (Timestamp can be different, then the restaurant itself).
func add_additional_fields(restaurant response_fields, id_from_reservation_page_url string, kitchen_office_hours restaurant_time) response_fields {
	restaurant_office_hours := get_opening_and_closing_time_from_restaurant_time(restaurant)

	restaurant.Links.TableReservationLocalizedId = id_from_reservation_page_url

	time_till_restaurant_closed := get_time_till_restaurant_closing_time(restaurant_office_hours.closing)
	time_till_restaurants_kitchen_closed := get_time_left_to_reserve(kitchen_office_hours.closing)

	restaurant.Openingtime.Time_till_restaurant_closed_hours = time_till_restaurant_closed.hour
	restaurant.Openingtime.Time_till_restaurant_closed_minutes = time_till_restaurant_closed.minutes

	restaurant.Openingtime.Time_left_to_reserve_hours = time_till_restaurants_kitchen_closed.hour
	restaurant.Openingtime.Time_left_to_reserve_minutes = time_till_restaurants_kitchen_closed.minutes

	restaurant.Links.TableReservationLocalizedId = id_from_reservation_page_url

	return restaurant
}

// Returns the id from a reservation page url associated with a restaurant because the id that comes with the restaurant might not be the same as the one in the
// reservation page url and that id is the one that lets us get access to the information related to reserving stuff etc.
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

func get_opening_and_closing_time_from_kitchen_time(restaurant response_fields) restaurant_time {
	// Converting restaurant_kitchen_start_time to unix, so we can compare it easily.
	// restaurant_kitchen_start_time := get_unix_from_time(restaurant.Openingtime.Restauranttime.Ranges[0].Start)
	restaurant_kitchen_start_time := get_unix_from_time(restaurant.Openingtime.Kitchentime.Ranges[0].Start)
	// We minus 1 hour from the end time because restaurants don't take reservations before that time slot.
	// IMPORTANT: E.g. if restaurant closes at 22:00, the last possible reservation time is 21:00.
	const one_hour_unix int64 = 3600
	// restaurant_kitchen_ending_time := get_unix_from_time(restaurant.Openingtime.Restauranttime.Ranges[0].End) - one_hour_unix
	restaurant_kitchen_ending_time := get_unix_from_time(restaurant.Openingtime.Kitchentime.Ranges[0].End) - one_hour_unix

	return restaurant_time{
		opening: restaurant_kitchen_start_time,
		closing: restaurant_kitchen_ending_time,
	}
}

func get_opening_and_closing_time_from_restaurant_time(restaurant response_fields) restaurant_time {
	// Converting restaurant_kitchen_start_time to unix, so we can compare it easily.
	// restaurant_kitchen_start_time := get_unix_from_time(restaurant.Openingtime.Restauranttime.Ranges[0].Start)
	restaurant_start_time := get_unix_from_time(restaurant.Openingtime.Restauranttime.Ranges[0].Start)
	// We minus 1 hour from the end time because restaurants don't take reservations before that time slot.
	// IMPORTANT: E.g. if restaurant closes at 22:00, the last possible reservation time is 21:00.
	const one_hour_unix int64 = 3600
	// restaurant_kitchen_ending_time := get_unix_from_time(restaurant.Openingtime.Restauranttime.Ranges[0].End) - one_hour_unix
	restaurant_ending_time := get_unix_from_time(restaurant.Openingtime.Restauranttime.Ranges[0].End) - one_hour_unix

	return restaurant_time{
		opening: restaurant_start_time,
		closing: restaurant_ending_time,
	}
}
func restaurant_format_is_incorrect(city string, restaurant response_fields) bool {

	// Converting to lower so that we don't run into problems when comparing it.
	city = strings.ToLower(city)

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

	restaurant_office_hours := get_opening_and_closing_time_from_kitchen_time(restaurant)
	// Checking to see if the timestamps are fucked here, so we don't have to check them later.
	// We have already checked that the ranges exist in the previous condition (restaurant.Openingtime.Restauranttime.Ranges != nil)
	return restaurant_office_hours.opening >= restaurant_office_hours.closing
}
