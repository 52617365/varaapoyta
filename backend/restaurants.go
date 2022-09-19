package main

import (
	"errors"
	"strings"

	"golang.org/x/exp/slices"
)

/*
This struct exists because it lets us filter out the restaurants we're not interested in (E.g. from city we didn't want)
whilst associating the response with the correct restaurant.
*/
type restaurant_with_time_slots struct {
	time_slots    chan parsed_graph_data
	restaurant    response_fields
	kitchen_times restaurant_time
}

// TODO: fixes currently being done.
func get_available_tables(city string, amount_of_eaters int) ([]response_fields, error) {
	current_time := get_current_date_and_time()
	all_time_intervals := get_all_raflaamo_time_intervals()
	time_slots_to_check := get_graph_time_slots_from_current_point_forward(current_time.time)

	raflaamo, err := init_restaurants()
	if err != nil {
		return []response_fields{}, errors.New("failed to construct the http client, contact the developer")
	}
	response, err := raflaamo.get()
	if err != nil {
		return []response_fields{}, errors.New("raflaamo api most likely down")
	}

	all_results := make(chan restaurant_with_time_slots, 50)
	for _, restaurant := range response {
		if restaurant_format_is_incorrect(city, restaurant) {
			continue
		}
		restaurant_additional_information := init_additional_information(restaurant, time_slots_to_check)

		is_err := restaurant_additional_information.add()
		if is_err != nil {
			// I.e. if the id was not found (invalid link)
			// This branch is most likely not taken.
			continue
		}
		restaurant_id := restaurant_additional_information.restaurant.Links.TableReservationLocalizedId
		kitchen_times := restaurant_additional_information.kitchen_office_hours

		// Demons start here.

		// Here we initialize this struct because we want to filter out all the unwanted restaurants so we don't have to iterate over a massive collection later (400 restaurants).
		// Also, here we get to associate the time slots that are not instantly known with the restaurant when the channel has the response ready.
		results := restaurant_with_time_slots{
			restaurant:    restaurant,
			time_slots:    make(chan parsed_graph_data, len(time_slots_to_check)),
			kitchen_times: kitchen_times,
		}
		all_results <- results
		graph_api_results := results.time_slots
		jobs := make(chan job, len(time_slots_to_check))

		// Spawning 8 workers.
		for i := 0; i < 8; i++ {
			go worker(jobs, graph_api_results)
		}

		// Spawning the jobs, this is 0-4 jobs every loop iteration on the response.
		for i := 0; i < len(time_slots_to_check); i++ {
			our_job := job{
				slot:             time_slots_to_check[i],
				restaurant_id:    restaurant_id,
				current_time:     current_time,
				amount_of_eaters: amount_of_eaters,
			}
			jobs <- our_job
		}
		close(jobs)
	}
	close(all_results)

	restaurants_with_opening_times := make([]response_fields, 50)
	for result := range all_results {
		restaurant, kitchen_times := result.restaurant, result.kitchen_times
		for time_slot := range result.time_slots {
			available_time_slots, err := extract_available_time_intervals_from_response(time_slot, current_time, kitchen_times, all_time_intervals)
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

type restaurant_time struct {
	opening int64
	closing int64
}

// Function contains several conditions which we check to determine if the restaurant is ok to use.
// We check them all inside of this one function.
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
