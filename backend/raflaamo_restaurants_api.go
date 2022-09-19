package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

type restaurants struct {
	request           *http.Request
	raflaamo_response chan []response_fields
	raflaamo_error    chan error
	http_client       *http.Client
}

// init_restaurants is used as a factory function to initialize the restaurants struct instance that is the used to call get()
func init_restaurants() (restaurants, error) {
	data := []byte(`{"operationName":"getRestaurantsByLocation","variables":{"first":1000,"input":{"restaurantType":"ALL","locationName":"Helsinki","feature":{"rentableVenues":false}},"after":"eyJmIjowLCJnIjp7ImEiOjYwLjE3MTE2LCJvIjoyNC45MzI1OH19"},"query":"fragment Locales on LocalizedString {fi_FI }fragment Restaurant on Restaurant {  id  name {    ...Locales    }  address {    municipality {      ...Locales       }        street {      ...Locales       }       zipCode     }   openingTime {    restaurantTime {      ranges {        start        end             }             }    kitchenTime {      ranges {        start        end        endNextDay              }             }    }  links {    tableReservationLocalized {      ...Locales        }    homepageLocalized {      ...Locales          }   }     }query getRestaurantsByLocation($first: Int, $after: String, $input: ListRestaurantsByLocationInput!) {  listRestaurantsByLocation(first: $first, after: $after, input: $input) {    totalCount      edges {      ...Restaurant        }     }}"}`)

	req, err := http.NewRequest("POST", "https://api.raflaamo.fi/query", bytes.NewBuffer(data))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("client_id", "jNAWMvWD9rp637RaR")
	req.Header.Add("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	if err != nil {
		return restaurants{}, errors.New("there was an error making the client in init_restaurants_api")
	}
	return restaurants{
		request:     req,
		http_client: &http.Client{},
	}, nil
}

// Sends a request to the raflaamo API and returns the deserialized response.
func (request *restaurants) get() ([]response_fields, error) {
	res, err := request.http_client.Do(request.request)

	if err != nil {
		request.raflaamo_error <- errors.New("there was an error connecting to the raflaamo api")
	}

	decoded, err := deserialize_response(&res)
	if err != nil {
		return []response_fields{}, errors.New("there was an error deserializing raflaamo API response")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	// Returning the start of the data, this will be an array.
	if len(decoded.Data.ListRestaurantsByLocation.Edges) == 0 {
		return []response_fields{}, errors.New("no restaurants found")
	}

	return decoded.Data.ListRestaurantsByLocation.Edges, nil
}
