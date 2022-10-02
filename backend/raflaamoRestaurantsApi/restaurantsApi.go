package raflaamoRestaurantsApi

import (
	"bytes"
	"errors"
	"log"
	"net/http"
)

type RaflaamoRestaurantsApi struct {
	httpClient *http.Client
	request    *http.Request
	response   *http.Response
}

func getRaflaamoRestaurantsApi() (*RaflaamoRestaurantsApi, error) {
	httpClient := &http.Client{}
	const data = `{"operationName":"getRestaurantsByLocation","variables":{"first":1000,"input":{"restaurantType":"ALL","locationName":"Helsinki","feature":{"rentableVenues":false}},"after":"eyJmIjowLCJnIjp7ImEiOjYwLjE3MTE2LCJvIjoyNC45MzI1OH19"},"query":"fragment Locales on LocalizedString {fi_FI }fragment Restaurant on Restaurant {  id  name {    ...Locales    }  address {    municipality {      ...Locales       }        street {      ...Locales       }       zipCode     }   openingTime {    restaurantTime {      ranges {        start        end             }             }    kitchenTime {      ranges {        start        end        endNextDay              }             }    }  links {    tableReservationLocalized {      ...Locales        }    homepageLocalized {      ...Locales          }   }     }query getRestaurantsByLocation($first: Int, $after: String, $input: ListRestaurantsByLocationInput!) {  listRestaurantsByLocation(first: $first, after: $after, input: $input) {    totalCount      edges {      ...Restaurant        }     }}"}`

	req, err := http.NewRequest("POST", "https://api.raflaamo.fi/query", bytes.NewBuffer([]byte(data)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("client_id", "jNAWMvWD9rp637RaR")
	req.Header.Add("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	if err != nil {
		return nil, errors.New("[getRaflaamoRestaurantsApi] error initializing RaflaamoRestaurantsApi")
	}

	return &RaflaamoRestaurantsApi{
		httpClient: httpClient,
		request:    req,
	}, nil
}

func noRestaurantsFound(restaurantsField *responseTopLevel) bool {
	return len(restaurantsField.Data.ListRestaurantsByLocation.Edges) == 0
}

func (restaurantsApi *RaflaamoRestaurantsApi) getRestaurantsFromRaflaamoApi() (*responseTopLevel, error) {
	restaurantsApi, err := getRaflaamoRestaurantsApi()
	if err != nil {
		log.Fatal(err)
	}

	httpClient := restaurantsApi.httpClient
	request := restaurantsApi.request

	res, err := httpClient.Do(request)

	if err != nil {
		return nil, errors.New("there was an error connecting to the raflaamo api")
	}
	restaurantsApi.response = res

	decodedRaflaamoRestaurants, err := restaurantsApi.deserializeRaflaamoRestaurantsResponse()
	if err != nil {
		return nil, errors.New("there was an error deserializing raflaamo API response")
	}

	if noRestaurantsFound(decodedRaflaamoRestaurants) {
		return nil, errors.New("no restaurants found with the specified criteria")
	}

	return decodedRaflaamoRestaurants, nil
}
