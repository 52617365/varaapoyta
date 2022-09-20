package main

import (
	"errors"
	"strings"
)

func get_available_tables(city string, amount_of_eaters int) ([]*response_fields, error) {
	current_time := get_current_date_and_time()
	all_time_intervals := get_all_raflaamo_time_intervals()
	time_slots_to_check := get_graph_time_slots_from_current_point_forward(current_time.time)

	raflaamo, err := init_restaurants()
	if err != nil {
		return nil, errors.New("failed to connect to the raflaamo api, contact the developer")
	}
	response, err := raflaamo.get()
	if err != nil {
		return nil, errors.New("raflaamo api most likely down")
	}

	all_results := make(chan additional_information, 50)
	for _, restaurant := range response {
		if restaurant_format_is_incorrect(city, &restaurant) {
			continue
		}
		restaurant_additional_information := init_additional_information(restaurant, len(time_slots_to_check))

		id_not_found_err := restaurant_additional_information.add()
		if id_not_found_err != nil {
			// This branch is most likely not taken.
			continue
		}
		restaurant_id := restaurant_additional_information.restaurant.Links.TableReservationLocalizedId

		// Storing the result into the slice of all results and modifying the result as a reference later.
		all_results <- restaurant_additional_information

		jobs := make(chan job, len(time_slots_to_check))

		// Spawning n workers, more than n is useless.
		for i := 0; i <= len(time_slots_to_check); i++ {
			go worker(jobs, restaurant_additional_information.time_slots)
		}

		// Spawning the jobs, this is 0-4 jobs every loop iteration on the response.
		for i := 0; i < len(time_slots_to_check); i++ {
			our_job := job{
				slot:             &time_slots_to_check[i],
				restaurant_id:    restaurant_id,
				current_time:     current_time,
				amount_of_eaters: amount_of_eaters,
			}
			jobs <- our_job
		}
		close(jobs)
	}

	restaurants_with_opening_times := make([]*response_fields, 0, 50)
	close(all_results)
	for result := range all_results {
		restaurant := result.restaurant
		kitchen_times := result.kitchen_times
		api_response := <-result.time_slots
		if api_response.err != nil {
			continue
		}
		for i := 0; i <= len(result.time_slots); i++ {
			available_time_slots, err := extract_available_time_intervals_from_response(api_response.value, current_time, &kitchen_times, all_time_intervals)
			if err != nil {
				continue
			}

			restaurant.Available_time_slots = available_time_slots
			restaurants_with_opening_times = append(restaurants_with_opening_times, &restaurant)
		}
	}
	// @Notice: If restaurant.Available_time_slots is null/nil, there are no available time slots.
	return restaurants_with_opening_times, nil
}

type restaurant_time struct {
	opening int64
	closing int64
}

// Function contains several conditions which we check to determine if the restaurant is ok to use.
// We check them all inside of this one function.
func restaurant_format_is_incorrect(city string, restaurant *response_fields) bool {
	// Converting to lower so that we don't run into problems when comparing it.
	city = strings.ToLower(city)

	// The restaurant times are nil if there are no time ranges available, we want to skip those because they are useless to us without times.
	if restaurant.Openingtime.Restauranttime.Ranges == nil || restaurant.Openingtime.Kitchentime.Ranges == nil {
		return true
	}
	if restaurant.Openingtime.Kitchentime.Ranges[0].Start == "" || restaurant.Openingtime.Kitchentime.Ranges[0].End == "" {
		return true
	}
	if !strings.Contains(restaurant.Links.TableReservationLocalized.Fi_FI, "https://s-varaukset.fi/online/reservation/fi/") {
		return true
	}
	// We check that the city is the correct one since we want to filter them out depending on the passed in city.
	if strings.ToLower(restaurant.Address.Municipality.Fi_FI) != city {
		return true
	}
	restaurant_office_hours := get_opening_and_closing_time_from_kitchen_time(restaurant)
	// Checking to see if the timestamps are fucked here, so we don't have to check them later.
	// We have already checked that the ranges exist in the previous condition (restaurant.Openingtime.Restauranttime.Ranges != nil)
	is_invalid_date := restaurant_office_hours.opening >= restaurant_office_hours.closing
	return is_invalid_date
}
