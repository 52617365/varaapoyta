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
	available_time_slots []time_slot_struct
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

func getAvailableTables(restaurants *[]response_fields, amount_of_eaters int) *[]restaurant_with_available_times_struct {
	re, _ := regexp.Compile(`[^fi/]\d+`) // This regex gets the first number match from the TableReservationLocalized JSON field which is the id we want. (https://regex101.com/r/NtFMrz/1)
	current_date := getCurrentDate()

	// Closest to a constant array we can get.
	// make a check to see if time is in the past, we don't care about the information if it's in the past. Might require we store times as integers and then convert to strings if needed.
	var all_possible_time_slots = [...]string{"0200", "0800", "1400", "2000"} // 02:00 covers(00:00-06:00), 08:00 covers(6:00-12:00), 14:00 covers(12:00-18:00), 20:00 covers(18:00-00:00)

	// there can be maximum of restaurants * all_possible_time_slots, so we allocate the worst case scenario here to avoid reallocation's.
	total_memory_to_reserve_for_all_restaurants := len(*restaurants) * len(all_possible_time_slots)
	all_restaurants_with_available_times := make([]restaurant_with_available_times_struct, 0, total_memory_to_reserve_for_all_restaurants)

	for _, restaurant := range *restaurants {
		if restaurant_does_not_contain_reservation_page(&restaurant) {
			continue
		}
		id_from_reservation_page_url := get_id_from_reservation_page_url(&restaurant, re)

		// If for some reason the reservation page url did not contain an id (regex returns empty string)
		if id_from_reservation_page_url == "" {
			continue
		}

		restaurant_with_available_times := restaurant_with_available_times_struct{
			restaurant:           restaurant,
			available_time_slots: make([]time_slot_struct, 0, len(all_possible_time_slots)),
		}

		for _, time_slot := range all_possible_time_slots {
			deserialized_graph_data, err := get_time_slots_from_graph_api(&id_from_reservation_page_url, current_date, &time_slot, amount_of_eaters)
			if err != nil {
				continue
			}
			if time_slot_does_not_contain_open_tables(deserialized_graph_data) {
				continue
			}

			// Here we have some kind of graph visible.
			unix_timestamp_struct_of_available_table := convert_unix_timestamp_to_finland(deserialized_graph_data)

			restaurant_with_available_times.available_time_slots = append(restaurant_with_available_times.available_time_slots, unix_timestamp_struct_of_available_table)
		}

		// TODO: don't let the times overlap. Use E.G. a hashmap to see if time already exists, if it does not, insert it.
		all_restaurants_with_available_times = append(all_restaurants_with_available_times, restaurant_with_available_times)
	}
	return &all_restaurants_with_available_times
}

// We do this because the id from the "Id" field is not always the same as the id needed in the reservation page.
func get_id_from_reservation_page_url(restaurant *response_fields, re *regexp.Regexp) string {
	reservation_page_url := *restaurant.Links.TableReservationLocalized.Fi_FI

	// there are some weird magic strings that will make regex fail so check that it's the link we're interested in.
	if reservation_page_url_is_not_valid(&reservation_page_url) {
		return ""
	}
	id_from_reservation_page_url := re.FindString(reservation_page_url)
	return id_from_reservation_page_url
}

func reservation_page_url_is_not_valid(reservation_page_url *string) bool {
	return !strings.Contains(*reservation_page_url, "https://s-varaukset.fi/online/reservation/fi")
}

func get_time_slots_from_graph_api(id_from_reservation_page_url *string, current_date *string, time_slot *string, amount_of_eaters int) (*parsed_graph_data, error) {
	// Example of an url to send the get request to.
	// https://s-varaukset.fi/api/recommendations/slot/{id}/{date}/{time}/{amount_of_eaters}
	request_url := fmt.Sprintf("https://s-varaukset.fi/api/recommendations/slot/%s/%s/%s/%d", *id_from_reservation_page_url, *current_date, *time_slot, amount_of_eaters)

	r, err := http.NewRequest("GET", request_url, nil)

	if err != nil {
		return nil, errors.New("error constructing get request")
	}
	r.Header.Add("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	client := &http.Client{}
	res, err := client.Do(r)

	if err != nil {
		return nil, errors.New("error sending get request")
	}

	deserialized_graph_data := deserialize_graph_response(&res)

	return deserialized_graph_data, nil
}

// We determine if there is a time slot with open tables by looking at the "color" field in the response.
// The color field will contain "transparent" if it does not contain a graph (open times), else it contains nil (there are open tables)
func time_slot_does_not_contain_open_tables(data *parsed_graph_data) bool {
	// "color" field is included and set to "transparent" if a graph does NOT exist on the page. (No times for restaurant).
	// Else it's nil.
	return data.Intervals[0].Color != nil
}

func restaurant_does_not_contain_reservation_page(restaurant *response_fields) bool {
	return len(*restaurant.Links.TableReservationLocalized.Fi_FI) == 0
}
