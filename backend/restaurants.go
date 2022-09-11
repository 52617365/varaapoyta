package main

import (
	"errors"
	"regexp"
	"strings"

	"github.com/gammazero/workerpool"
)

// TODO: Even after our recent changes, this function is unfortunately still too slow (10 seconds for helsinki restaurants), we NEED to optimize it. Now is the time.
func get_available_tables(city string, amount_of_eaters int) []response_fields {

	// Using channel because this function would block for around ~2 seconds otherwise.
	raflaamo_api_response := make(chan []response_fields)
	raflaamo_api_response_error := make(chan error)
	go get_all_restaurants_from_raflaamo_api(raflaamo_api_response, raflaamo_api_response_error)

	// Getting current_time, so we can avoid checking times from the past.
	current_time := get_current_date_and_time()
	// All possible time slots we need to check, it does not contain time slots from the past.
	time_slots_to_check_from_graph_api := get_graph_time_slots_from_current_point_forward(current_time.time)
	// 11:00, 11:15, 11:30 and so on.
	all_time_intervals := get_all_raflaamo_time_intervals()

	// Checking if there was an error when requesting the raflaamo API.
	if <-raflaamo_api_response_error != nil {
		return nil
	}
	// Getting the response received from the raflaamo api.
	// We capture it here to avoid blocking for the duration of the request.
	// Here we have already checked that it wasn't an error and can treat it as valid.
	raflaamo_api_restaurants := <-raflaamo_api_response

	// Channel that will hold all the restaurants with additional information related to their opening times.
	restaurants_chan := make(chan response_fields, len(raflaamo_api_restaurants))

	wp := workerpool.New(8)
	for _, restaurant := range raflaamo_api_restaurants {
		restaurant := restaurant
		wp.Submit(func() {
			if restaurant_format_is_incorrect(city, restaurant) {
				return
			}
			// If we can't find the id from url, just return on to the next one because without the id we can't find the reservation page.
			id_from_reservation_page_url, err := get_id_from_reservation_page_url(restaurant)
			if err != nil {
				// Not finding a valid id results in an err, without a valid id we can't check any restaurant related to the id, so we just return.
				return
			}

			kitchen_office_hours := get_opening_and_closing_time_from_kitchen_time(restaurant)
			restaurant_office_hours := get_opening_and_closing_time_from_restaurant_time(restaurant)
			available_intervals_from_graph_api, err := get_available_time_intervals_from_graph_api(kitchen_office_hours.opening, kitchen_office_hours.closing, id_from_reservation_page_url, time_slots_to_check_from_graph_api, amount_of_eaters, all_time_intervals, current_time)
			if err != nil {
				return
			}
			// Here we populate the empty field time slot with all the available time slots.
			// This is expected behavior because we planned on populating it later on.

			restaurant.Available_time_slots = available_intervals_from_graph_api
			// Storing this, so we can easily use it later without parsing the reservation url again.
			restaurant.Links.TableReservationLocalizedId = id_from_reservation_page_url

			// Adding relative times, so we can display them as a countdown on the page later.
			time_till_restaurant_closed := get_time_till_restaurant_closing_time(restaurant_office_hours.closing)
			time_till_restaurants_kitchen_closed := get_time_till_restaurant_closing_time(kitchen_office_hours.closing)

			// Adding the relative time to when the restaurant itself closes (Timestamp can be different, then the kitchen time).
			restaurant.Openingtime.Time_till_restaurant_closed_hours = time_till_restaurant_closed.hour
			restaurant.Openingtime.Time_till_restaurant_closed_minutes = time_till_restaurant_closed.minutes

			// Adding the relative time to when the restaurants kitchen closes (Timestamp can be different, then the restaurant itself).
			restaurant.Openingtime.Time_till_kitchen_closed_hours = time_till_restaurants_kitchen_closed.hour
			restaurant.Openingtime.Time_till_kitchen_closed_minutes = time_till_restaurants_kitchen_closed.minutes

			restaurant.Links.TableReservationLocalizedId = id_from_reservation_page_url
			restaurants_chan <- restaurant
		})
	}
	wp.StopWait()
	close(restaurants_chan)
	// Appending channel into slice instead of directly returning channel because we serialize the array afterwards.
	restaurants_from_provided_city := make([]response_fields, 0, len(raflaamo_api_restaurants))
	for restaurant := range restaurants_chan {
		restaurants_from_provided_city = append(restaurants_from_provided_city, restaurant)
	}
	return restaurants_from_provided_city
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
	restaurant_kitchen_start_time := get_unix_from_time(restaurant.Openingtime.Restauranttime.Ranges[0].Start)
	// We minus 1 hour from the end time because restaurants don't take reservations before that time slot.
	// IMPORTANT: E.g. if restaurant closes at 22:00, the last possible reservation time is 21:00.
	const one_hour_unix int64 = 3600
	// restaurant_kitchen_ending_time := get_unix_from_time(restaurant.Openingtime.Restauranttime.Ranges[0].End) - one_hour_unix
	restaurant_kitchen_ending_time := get_unix_from_time(restaurant.Openingtime.Restauranttime.Ranges[0].End) - one_hour_unix

	return restaurant_time{
		opening: restaurant_kitchen_start_time,
		closing: restaurant_kitchen_ending_time,
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
