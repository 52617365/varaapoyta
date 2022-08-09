package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
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

// generatePayloadsFromIdAndSend gets all possible payloads from a certain time passed in as an argument and sends all of them.
// E.g if function is called at 13:00 it generates payloads from 13:30 onwards.
// @ Use worker pool with goroutines?
// Returns an array of all available times as string array.
func generatePayloadsFromIdAndSend(id *string) (*[]string, error) {
	times := getAllPossibleTimes() // all times from current time forward.
	current_date := getCurrentDate()

	available_tables_from_id := make([]string, 0, len(*times))

	for _, time := range *times {
		// Making the payload we later send get request to.
		url_formatted_time := strings.Replace(time, ":", "%3A", -1) // Replace ":" with browser equivalent "%3A"
		payload := fmt.Sprintf("https://s-varaukset.fi/online/reserve/availability/fi/%s?date=%s&slot_id=%s&time=%s&amount=1&price_code=&check=1", *id, *current_date, *id, url_formatted_time)

		r, err := http.NewRequest("GET", payload, nil)
		r.Header.Add("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

		if err != nil {
			return &available_tables_from_id, errors.New("error initializing http client")
		}

		client := &http.Client{}
		res, err := client.Do(r)

		if err != nil {
			return &available_tables_from_id, errors.New("error sending request")
		}

		resBody, err := io.ReadAll(res.Body)

		if err != nil {
			return &available_tables_from_id, errors.New("error reading response body")
		}

		// Checking if the table can be reserved and if it can, appending to reservable tables.
		if strings.Contains(string(resBody), "Varauksen tekeminen hakemallenne ajankohdalle on mahdollista") {
			available_tables_from_id = append(available_tables_from_id, time)
		}
	}
	return &available_tables_from_id, nil
}

func getAvailableTables(restaurants *[]response_fields) *[]available_times {
	available_tables := make([]available_times, 0, len(*restaurants))

	// https://regex101.com/r/NtFMrz/1
	// This regex gets the first number match from the TableReservationLocalized JSON field.
	re := regexp.MustCompile(`[^fi/]\d+`)

	for _, restaurant := range *restaurants {
		if restaurant_does_not_contain_reservation_page(&restaurant) {
			continue
		}

		reservation_page_url := *restaurant.Links.TableReservationLocalized.Fi_FI
		id := re.FindString(reservation_page_url)
		// @ Use goroutines and channels?
		available_tables_from_id, err := generatePayloadsFromIdAndSend(&id)
		if id_does_not_contain_open_tables(err, available_tables_from_id) {
			continue
		}

		// Here we can assume that there are available times.
		struct_from_available_tables := available_times{
			restaurant: &restaurant,
			times:      available_tables_from_id,
		}
		available_tables = append(available_tables, struct_from_available_tables)
	}
	return &available_tables
}

func id_does_not_contain_open_tables(err error, available_tables_from_id *[]string) bool {
	return (err != nil || len(*available_tables_from_id) == 0)
}

func restaurant_does_not_contain_reservation_page(restaurant *response_fields) bool {
	return (len(*restaurant.Links.TableReservationLocalized.Fi_FI) == 0)
}
