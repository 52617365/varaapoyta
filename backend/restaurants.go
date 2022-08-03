package main

import (
	"bytes"
	"log"
	"net/http"
)

// TODO: pass city here.
func getRestaurants() []response_fields {
	data := []byte(`{"operationName":"getRestaurantsByLocation","variables":{"first":470,"input":{"restaurantType":"ALL","locationName":"Helsinki","feature":{"rentableVenues":false}},"after":"eyJmIjoxMCwiZyI6eyJhIjo2MC4xNzExNiwibyI6MjQuOTMyNTh9fQo"},"query":"fragment Locales on LocalizedString {fi_FI }fragment Restaurant on Restaurant {  id  name {    ...Locales    }  urlPath {    ...Locales     }    address {    municipality {      ...Locales       }        street {      ...Locales       }       zipCode     }    features {    accessible     }  openingTime {    restaurantTime {      ranges {        start        end        endNextDay         }             }    kitchenTime {      ranges {        start        end        endNextDay              }             }    }  links {    tableReservationLocalized {      ...Locales        }    homepageLocalized {      ...Locales          }   }     }query getRestaurantsByLocation($first: Int, $after: String, $input: ListRestaurantsByLocationInput!) {  listRestaurantsByLocation(first: $first, after: $after, input: $input) {    totalCount      edges {      ...Restaurant        }     }}"}`)

	r, err := http.NewRequest("POST", "https://api.raflaamo.fi/query", bytes.NewBuffer(data))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("client_id", "jNAWMvWD9rp637RaR")

	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	res, err := client.Do(r)

	if err != nil {
		log.Fatal(err)
	}

	//response, _ := io.ReadAll(res.Body)
	//fmt.Println(string(response))
	decoded := deserialize_response(&res)
	defer res.Body.Close()

	// Returning the start of the data, this will be an array.
	// Later we just loop over all of this data and do w/e.
	return decoded.Data.ListRestaurantsByLocation.Edges
}
