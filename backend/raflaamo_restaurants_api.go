package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

// Sends a request to the raflaamo API and returns the deserialized response.
// TODO: when searching for helsinki, we get an index out of range error. Fix this.
func get_all_restaurants_from_raflaamo_api(result chan<- []response_fields, possible_error chan<- error) {
	data := []byte(`{"operationName":"getRestaurantsByLocation","variables":{"first":1000,"input":{"restaurantType":"ALL","locationName":"Helsinki","feature":{"rentableVenues":false}},"after":"eyJmIjowLCJnIjp7ImEiOjYwLjE3MTE2LCJvIjoyNC45MzI1OH19"},"query":"fragment Locales on LocalizedString {fi_FI }fragment Restaurant on Restaurant {  id  name {    ...Locales    }  address {    municipality {      ...Locales       }        street {      ...Locales       }       zipCode     }   openingTime {    restaurantTime {      ranges {        start        end             }             }    kitchenTime {      ranges {        start        end        endNextDay              }             }    }  links {    tableReservationLocalized {      ...Locales        }    homepageLocalized {      ...Locales          }   }     }query getRestaurantsByLocation($first: Int, $after: String, $input: ListRestaurantsByLocationInput!) {  listRestaurantsByLocation(first: $first, after: $after, input: $input) {    totalCount      edges {      ...Restaurant        }     }}"}`)

	r, err := http.NewRequest("POST", "https://api.raflaamo.fi/query", bytes.NewBuffer(data))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("client_id", "jNAWMvWD9rp637RaR")
	r.Header.Add("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	if err != nil {
		possible_error <- errors.New("there was an error connecting to the raflaamo api")
	}
	client := &http.Client{}
	res, err := client.Do(r)

	if err != nil {
		possible_error <- errors.New("there was an error connecting to the raflaamo api")
	}

	decoded, err := deserialize_response(&res)
	if err != nil {
		possible_error <- errors.New("there was an error deserializing raflaamo API response")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	// Returning the start of the data, this will be an array.
	if len(decoded.Data.ListRestaurantsByLocation.Edges) == 0 {
		possible_error <- errors.New("no restaurants found")
	}

	possible_error <- nil
	result <- decoded.Data.ListRestaurantsByLocation.Edges
	// Closing channel because this is the only sender to this specific channel.
	close(result)
	close(possible_error)
}
