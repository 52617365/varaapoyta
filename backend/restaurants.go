package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

// Contains the restaurant information and on top of that, all available times you can reserve a table from that restaurant.
type available_times struct {
	restaurant *response_fields
	times      *[]string
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

func coverAllAvailableTimeSlots(id *string, amount_of_eaters int) {
	current_date := getCurrentDate()
	// Time(Covered time) One time covers a 6h time window in graph.
	// 02:00(00:00-06:00), 08:00(6:00-12:00), 14:00(12:00-18:00), 20:00(18:00-00:00)
	all_time_slots := []string{"0200", "0800", "1400", "2000"}
	time_slot_results := make([]parsed_graph_data, 0, len(all_time_slots))

	for _, time_slot := range all_time_slots {
		// https://s-varaukset.fi/api/recommendations/slot/{id}/{date}/{time}/{amount_of_eaters}
		payload := fmt.Sprintf("https://s-varaukset.fi/api/recommendations/slot/%s/%s/%s/%d", *id, *current_date, time_slot, amount_of_eaters)

		r, err := http.NewRequest("GET", payload, nil)
		r.Header.Add("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

		if err != nil {
			continue
		}

		client := &http.Client{}
		res, err := client.Do(r)

		if err != nil {
			continue
		}
		parsed_graph_data := deserialize_graph_response(&res)
		fmt.Println(parsed_graph_data.Id)
		time_slot_results = append(time_slot_results, *parsed_graph_data)
	}
}

func getAvailableTables(restaurants *[]response_fields) *[]available_times {
	available_tables := make([]available_times, 0, len(*restaurants))

	// https://regex101.com/r/NtFMrz/1
	// This regex gets the first number match from the TableReservationLocalized JSON field which is the id we want.
	re := regexp.MustCompile(`[^fi/]\d+`)

	for _, restaurant := range *restaurants {
		if restaurant_does_not_contain_reservation_page(&restaurant) {
			continue
		}

		// TODO: use goroutines, its hella slow rn.
		reservation_page_url := *restaurant.Links.TableReservationLocalized.Fi_FI
		id := re.FindString(reservation_page_url)
		// TODO: capture available table slots from the return value of this function call.
		coverAllAvailableTimeSlots(&id, 1)

		// TODO: modify this struct into the new coverAllAvailableTimeSlots function result.
		// Here we can assume that there are available times.
		// struct_from_available_tables := available_times{
		// 	restaurant: &restaurant,
		// 	times:      available_tables_from_id,
		// }
		// available_tables = append(available_tables, struct_from_available_tables)
	}
	return &available_tables
}

// mby use this again but with a different check (check "color" field if it exists or something)
func id_does_not_contain_open_tables(err error, available_tables_from_id *[]string) bool {
	return err != nil || len(*available_tables_from_id) == 0
}

func restaurant_does_not_contain_reservation_page(restaurant *response_fields) bool {
	return len(*restaurant.Links.TableReservationLocalized.Fi_FI) == 0
}
