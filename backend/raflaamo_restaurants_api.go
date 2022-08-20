package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

// Contains the restaurant information and on top of that, all available times you can reserve a table from that restaurant.
type restaurant_with_available_times_struct struct {
	restaurant           response_fields
	available_time_slots []string
}

// Sends a request to the raflaamo API and returns the deserialized response. err != nil if failed.
func get_all_restaurants_from_raflaamo_api() ([]response_fields, error) {
	data := []byte(`{"operationName":"getRestaurantsByLocation","variables":{"first":1000,"input":{"restaurantType":"ALL","locationName":"Helsinki","feature":{"rentableVenues":false}},"after":"eyJmIjowLCJnIjp7ImEiOjYwLjE3MTE2LCJvIjoyNC45MzI1OH19"},"query":"fragment Locales on LocalizedString {fi_FI }fragment Restaurant on Restaurant {  id  name {    ...Locales    }  urlPath {    ...Locales     }    address {    municipality {      ...Locales       }        street {      ...Locales       }       zipCode     }    features {    accessible     }  openingTime {    restaurantTime {      ranges {        start        end             }             }    kitchenTime {      ranges {        start        end        endNextDay              }             }    }  links {    tableReservationLocalized {      ...Locales        }    homepageLocalized {      ...Locales          }   }     }query getRestaurantsByLocation($first: Int, $after: String, $input: ListRestaurantsByLocationInput!) {  listRestaurantsByLocation(first: $first, after: $after, input: $input) {    totalCount      edges {      ...Restaurant        }     }}"}`)

	r, err := http.NewRequest("POST", "https://api.raflaamo.fi/query", bytes.NewBuffer(data))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("client_id", "jNAWMvWD9rp637RaR")
	r.Header.Add("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	if err != nil {
		return nil, errors.New("there was an error connecting to the raflaamo api")
	}
	client := &http.Client{}
	res, err := client.Do(r)

	if err != nil {
		return nil, errors.New("there was an error connecting to the raflaamo api")
	}

	decoded, err := deserialize_response(&res)
	if err != nil {
		return nil, errors.New("there was an error deserializing raflaamo API response")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	// Returning the start of the data, this will be an array.
	return decoded.Data.ListRestaurantsByLocation.Edges, nil
}
