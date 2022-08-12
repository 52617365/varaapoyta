package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

type time_slot struct {
	start *string
	end   *string
}

// Contains the restaurant information and on top of that, all available times you can reserve a table from that restaurant.
type available_times struct {
	restaurant          *response_fields
	available_time_slot *time_slot
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

func getAvailableTables(restaurants *[]response_fields) *[]available_times {
	available_tables := make([]available_times, 0, len(*restaurants))
	re := regexp.MustCompile(`[^fi/]\d+`) // This regex gets the first number match from the TableReservationLocalized JSON field which is the id we want. (https://regex101.com/r/NtFMrz/1)
	current_date := getCurrentDate()

	// Closest to a constant array we can get.
	var all_possible_time_slots = [...]string{"0200", "0800", "1400", "2000"} // 02:00 covers(00:00-06:00), 08:00 covers(6:00-12:00), 14:00 covers(12:00-18:00), 20:00 covers(18:00-00:00)

	for _, restaurant := range *restaurants {
		if restaurant_does_not_contain_reservation_page(&restaurant) {
			continue
		}

		reservation_page_url := *restaurant.Links.TableReservationLocalized.Fi_FI
		id_from_reservation_page_url := re.FindString(reservation_page_url)

		available_times_from_restaurant := make([]parsed_graph_data, 0, len(all_possible_time_slots))

		// TODO: use goroutines.
		for _, time_slot := range all_possible_time_slots {
			// https://s-varaukset.fi/api/recommendations/slot/{id}/{date}/{time}/{amount_of_eaters}
			request_url := fmt.Sprintf("https://s-varaukset.fi/api/recommendations/slot/%s/%s/%s/%d", id_from_reservation_page_url, *current_date, time_slot, 1)

			r, err := http.NewRequest("GET", request_url, nil)
			r.Header.Add("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

			if err != nil {
				continue
			}

			client := &http.Client{}
			res, err := client.Do(r)

			if err != nil {
				continue
			}
			deserialized_graph_data := deserialize_graph_response(&res)
			if time_slot_does_not_contain_open_tables(deserialized_graph_data) {
				continue
			}
			// Here we have some kind of graph (times available for the time_slot)
			// TODO: figure out how we want to store the available table data.
			available_times_from_restaurant = append(available_times_from_restaurant, *deserialized_graph_data)
		}
	}
	return &available_tables
}

// mby use this again but with a different check (check "color" field if it exists or something)
func time_slot_does_not_contain_open_tables(data *parsed_graph_data) bool {
	// "color" field is included and set to "transparent" if a graph does NOT exist on the page. (No times for restaurant).
	// Else it's nil.
	if data.Intervals[0].Color != nil { // Here it's transparent aka there are no free tables.
		return true
	}
	return false
}

func restaurant_does_not_contain_reservation_page(restaurant *response_fields) bool {
	return len(*restaurant.Links.TableReservationLocalized.Fi_FI) == 0
}
