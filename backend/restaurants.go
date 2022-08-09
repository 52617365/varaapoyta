package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func getRestaurants() *[]response_fields {
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

func getPayloads(id *string) *[]string {
	times := getAllPossibleTimes() // all times from current time forward.
	payload_strings := make([]string, 0, len(*times))
	current_date := getCurrentDate()

	// Example string: https://s-varaukset.fi/online/reserve/availability/fi/38?date=2022-08-24&slot_id=38&time=12%3A15&amount=1&price_code=&check=1
	for _, time := range *times {
		url_formatted_time := strings.Replace(time, ":", "%3A", -1)
		payload := fmt.Sprintf("https://s-varaukset.fi/online/reserve/availability/fi/%s?date=%s&slot_id=%s&time=%s&amount=1&price_code=&check=1", id, current_date, id, url_formatted_time)
		payload_strings = append(payload_strings, payload)
	}
	return &payload_strings
}

func getAvailableTables(restaurants *[]response_fields) {
	for _, restaurant := range *restaurants {
		if len(*restaurant.Links.TableReservationLocalized.Fi_FI) != 0 {
			id := *&restaurant.Id
			payloads := getPayloads(id)
			// TODO: do something with payload instead of printing it.
			fmt.Println(payloads)
		}

	}
}
