package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// Contains the restaurant information and on top of that, all available times you can reserve a table from that restaurant.
type restaurant_with_available_times_struct struct {
	restaurant           response_fields
	available_time_slots []string
}

// getAllRestaurantsFromRaflaamoApi sends request to raflaamo API and gets all the possible restaurants from all cities from there.
func getAllRestaurantsFromRaflaamoApi() *[]response_fields {
	data := []byte(`{"operationName":"getRestaurantsByLocation","variables":{"first":1000,"input":{"restaurantType":"ALL","locationName":"Helsinki","feature":{"rentableVenues":false}},"after":"eyJmIjowLCJnIjp7ImEiOjYwLjE3MTE2LCJvIjoyNC45MzI1OH19"},"query":"fragment Locales on LocalizedString {fi_FI }fragment Restaurant on Restaurant {  id  name {    ...Locales    }  urlPath {    ...Locales     }    address {    municipality {      ...Locales       }        street {      ...Locales       }       zipCode     }    features {    accessible     }  openingTime {    restaurantTime {      ranges {        start        end             }             }    kitchenTime {      ranges {        start        end        endNextDay              }             }    }  links {    tableReservationLocalized {      ...Locales        }    homepageLocalized {      ...Locales          }   }     }query getRestaurantsByLocation($first: Int, $after: String, $input: ListRestaurantsByLocationInput!) {  listRestaurantsByLocation(first: $first, after: $after, input: $input) {    totalCount      edges {      ...Restaurant        }     }}"}`)

	r, err := http.NewRequest("POST", "https://api.raflaamo.fi/query", bytes.NewBuffer(data))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("client_id", "jNAWMvWD9rp637RaR")
	r.Header.Add("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	res, err := client.Do(r)

	if err != nil {
		log.Fatal(err)
	}

	decoded := deserialize_response(&res)
	defer res.Body.Close()

	// Returning the start of the data, this will be an array.
	return decoded.Data.ListRestaurantsByLocation.Edges
}

// TODO: use goroutines for requests
func getAvailableTables(restaurants *[]response_fields, amount_of_eaters int) *[]restaurant_with_available_times_struct {
	re, err := regexp.Compile(`[^fi/]\d+`) // This regex gets the first number match from the TableReservationLocalized JSON field which is the id we want. https://regex102.com/r/NtFMrz/1
	if err != nil {
		log.Fatal("why the fuck did this fail?")
	}
	current_date := getCurrentDate()

	// TODO: make a check to see if time is in the past, we don't care about the information if it's in the past. (get_current_time)
	var all_possible_time_slots = [...]string{"0200", "0800", "1400", "2000"} // 02:00 covers(00:00-06:00), 08:00 covers(6:00-12:00), 14:00 covers(12:00-18:00), 20:00 covers(18:00-00:00)

	// There can be maximum of restaurants * all_possible_time_slots, so we allocate the worst case scenario here to avoid reallocation's.
	total_memory_to_reserve_for_all_restaurant_time_slots := len(*restaurants) * len(all_possible_time_slots)

	// This will contain all the available time slots from all restaurants after loop runs.
	all_restaurants_with_available_times := make([]restaurant_with_available_times_struct, 0, total_memory_to_reserve_for_all_restaurant_time_slots)

	for _, restaurant := range *restaurants {
		id_from_reservation_page_url, err := get_id_from_reservation_page_url(&restaurant, re)

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
			time_slots_from_graph_api, err := get_time_slots_from_graph_api(id_from_reservation_page_url, *current_date, time_slot, amount_of_eaters)
			if err != nil {
				continue
			}

			// At this point in the code we have already made all the necessary checks to confirm that a graph is visible for a time slot, and we can extract information from it.

			unix_timestamp_struct_of_available_table := convert_unix_timestamp_to_finland_time(time_slots_from_graph_api)

			time_slots_in_between_unix_timestamps, err := return_time_slots_in_between(unix_timestamp_struct_of_available_table.start_time, unix_timestamp_struct_of_available_table.end_time)

			if err != nil {
				continue
			}

			// Checks to see if single_restaurant_with_available_times does not already contain the value we're about to add then adds it.
			add_non_duplicate_time_slots_into_array(time_slots_in_between_unix_timestamps, single_restaurant_with_available_times)
		}
		// Here after iterating over all time slots for the restaurant, we store the results.
		all_restaurants_with_available_times = append(all_restaurants_with_available_times, single_restaurant_with_available_times)
	}
	return &all_restaurants_with_available_times
}

// Iterates over all time slots and adds the iterated time slot into an array of time slots if the array does not already
// contain that specific time slot. This is done to only store unique time slots instead of having duplicates.
func add_non_duplicate_time_slots_into_array(time_slots_in_between_unix_timestamps *[]string, restaurant_with_available_times restaurant_with_available_times_struct) {
	for _, available_time_slot := range *time_slots_in_between_unix_timestamps {
		if !contains(&restaurant_with_available_times.available_time_slots, available_time_slot) {
			restaurant_with_available_times.available_time_slots = append(restaurant_with_available_times.available_time_slots, available_time_slot)
		}
	}
}

// Basic function to see if a container contains an element.
func contains(arr *[]string, our_string string) bool {
	for _, v := range *arr {
		if v == our_string {
			return true
		}
	}
	return false
}

// We do this because the id from the "Id" field is not always the same as the id needed in the reservation page.
func get_id_from_reservation_page_url(restaurant *response_fields, re *regexp.Regexp) (string, error) {
	reservation_page_url := restaurant.Links.TableReservationLocalized.Fi_FI
	if restaurant_does_not_contain_reservation_page(restaurant) {
		return "", errors.New("restaurant did not contain reservation page url")
	}
	// There are some weird magic strings that will make regex fail so check that it's the link we're interested in.
	if reservation_page_url_is_not_valid(&reservation_page_url) {
		return "", errors.New("reservation_page_url_is_not_valid")
	}
	id_from_reservation_page_url := re.FindString(reservation_page_url)

	// If regex could not match or if url was invalid (happens sometimes cuz API is weird).
	if id_from_reservation_page_url == "" {
		return "", errors.New("regex did not match anything, something wrong with reservation_page_url")
	}
	return id_from_reservation_page_url, nil
}

// Gets timeslots from raflaamo API that is responsible for returning graph data.
// Instead of drawing a graph with it, we convert it into time to determine which table is open or not.

func get_time_slots_from_graph_api(id_from_reservation_page_url string, current_date string, time_slot string, amount_of_eaters int) (*parsed_graph_data, error) {

	// https://s-varaukset.fi/api/recommendations/slot/{id}/{date}/{time}/{amount_of_eaters}
	request_url := fmt.Sprintf("https://s-varaukset.fi/api/recommendations/slot/%s/%s/%s/%d", id_from_reservation_page_url, current_date, time_slot, amount_of_eaters)

	r, err := http.NewRequest("GET", request_url, nil)

	if err != nil {
		return nil, errors.New("error constructing get request")
	}
	r.Header.Add("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	client := &http.Client{}
	res, err := client.Do(r)

	// Will throw if we call deserialize_graph_response with a status code other than 200 so we handle it here.
	if err != nil || res.StatusCode != 200 {
		return nil, errors.New("error sending get request")
	}

	deserialized_graph_data := deserialize_graph_response(&res)

	// most likely wont jump into this branch but check regardless.
	if deserialized_graph_data == nil {
		return nil, errors.New("there was an error deserializing the data returned from endpoint")
	}

	if time_slot_does_not_contain_open_tables(deserialized_graph_data) {
		return nil, errors.New("time slot did not contain open tables")
	}

	return deserialized_graph_data, nil
}

// Checks to see if reservation_page_url contains the correct url, sometimes the url is something related to renting a table
// Which will result in an invalid regex match when trying to get id from reservation_page_url.
// @This is raflaamo's fault, but we have to deal with it.
func reservation_page_url_is_not_valid(reservation_page_url *string) bool {
	return !strings.Contains(*reservation_page_url, "https://s-varaukset.fi/online/reservation/fi")
}

// We determine if there is a time slot with open tables by looking at the "color" field in the response.
// The color field will contain "transparent" if it does not contain a graph (open times), else it contains nil (meaning there are open tables)
func time_slot_does_not_contain_open_tables(data *parsed_graph_data) bool {
	return data.Intervals[0].Color != nil
}

// Some restaurants don't even contain a reservation page url, these restaurants are useless to us so we make sure to check it.
func restaurant_does_not_contain_reservation_page(restaurant *response_fields) bool {
	return len(restaurant.Links.TableReservationLocalized.Fi_FI) == 0
}
