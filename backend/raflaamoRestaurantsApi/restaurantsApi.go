/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package raflaamoRestaurantsApi

import (
	"backend/helpers"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func GetRaflaamoRestaurantsApi(city string) *RaflaamoRestaurantsApi {
	httpClient := &http.Client{}
	const data = `{"operationName":"getRestaurantsByLocation","variables":{"first":1000,"input":{"restaurantType":"ALL","locationName":"Helsinki","feature":{"rentableVenues":false}},"after":"eyJmIjowLCJnIjp7ImEiOjYwLjE3MTE2LCJvIjoyNC45MzI1OH19"},"query":"fragment Locales on LocalizedString {fi_FI }fragment Restaurant on Restaurant {  id  name {    ...Locales    }  address {    municipality {      ...Locales       }        street {      ...Locales       }       zipCode     }   openingTime {    restaurantTime {      ranges {        start        end             }             }    kitchenTime {      ranges {        start        end        endNextDay              }             }    }  links {    tableReservationLocalized {      ...Locales        }    homepageLocalized {      ...Locales          }   }     }query getRestaurantsByLocation($first: Int, $after: String, $input: ListRestaurantsByLocationInput!) {  listRestaurantsByLocation(first: $first, after: $after, input: $input) {    totalCount      edges {      ...Restaurant        }     }}"}`

	req, err := http.NewRequest("POST", "https://api.raflaamo.fi/query", bytes.NewBuffer([]byte(data)))

	if err != nil {
		log.Fatalln("[GetRaflaamoRestaurantsApi] - err but shouldn't be")
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("client_id", "jNAWMvWD9rp637RaR")
	req.Header.Add("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	return &RaflaamoRestaurantsApi{
		httpClient:               httpClient,
		request:                  req,
		cityToGetRestaurantsFrom: city,
	}
}

func (raflaamoRestaurantsApi *RaflaamoRestaurantsApi) GetRestaurantsFromRaflaamoApi(currentTime int64) ([]ResponseFields, error) {
	httpClient := raflaamoRestaurantsApi.httpClient
	request := raflaamoRestaurantsApi.request

	raflaamoRestaurantsApi.currentTime = currentTime

	res, err := httpClient.Do(request)

	if err != nil {
		return nil, fmt.Errorf("[GetRestaurantsFromRaflaamoApi] - %w", errors.New("could not get raflaamo restaurants since the site seems to be down"))
	}

	decodedRaflaamoRestaurants, err := raflaamoRestaurantsApi.deserializeRaflaamoRestaurantsResponse(res)

	if err != nil {
		return nil, fmt.Errorf("[GetRestaurantsFromRaflaamoApi] - %w", errors.New("there was an error deserializing raflaamo restaurants api response"))
	}

	filteredGoodRestaurants := raflaamoRestaurantsApi.filterBadRestaurantsOut(decodedRaflaamoRestaurants)

	return filteredGoodRestaurants, nil
}

// A restaurant is considered "Bad" if:
//   - Restaurants city is not from the provided city.
//   - Restaurants reservation link does not exist or contains odd contents.
//   - Restaurant does not contain opening times (Specified in the Ranges array).
//   - Restaurant or restaurants kitchen is already closed.
func (raflaamoRestaurantsApi *RaflaamoRestaurantsApi) filterBadRestaurantsOut(structureContainingRestaurantData *responseTopLevel) []ResponseFields {
	restaurantsFromRaflaamoApi := structureContainingRestaurantData.Data.ListRestaurantsByLocation.Edges

	filteredGoodRestaurants := make([]ResponseFields, 0, 40)
	for _, restaurant := range restaurantsFromRaflaamoApi {
		if restaurant.isBad(raflaamoRestaurantsApi.cityToGetRestaurantsFrom, raflaamoRestaurantsApi.currentTime) {
			continue
		}
		filteredGoodRestaurants = append(filteredGoodRestaurants, restaurant)
	}
	return filteredGoodRestaurants
}

func (response *ResponseFields) isBad(city string, currentTime int64) bool {
	if response.cityDoesNotMatchUsersCity(city) {
		return true
	}
	if response.reservationLinkIsNotValid() {
		return true
	}
	if response.doesNotContainOpeningTimes() {
		return true
	}
	if response.restaurantOrKitchenIsAlreadyClosed(currentTime) {
		return true
	}
	return false
}

func (response *ResponseFields) doesNotContainOpeningTimes() bool {
	restaurantsOpeningTimes := response.Openingtime.Restauranttime.Ranges
	kitchensOpeningTimes := response.Openingtime.Kitchentime.Ranges

	if restaurantsOpeningTimes == nil || kitchensOpeningTimes == nil {
		return true
	}
	if restaurantsOpeningTimes[0].Start == "" || restaurantsOpeningTimes[0].End == "" {
		return true
	}
	return false
}

func (response *ResponseFields) restaurantOrKitchenIsAlreadyClosed(currentTime int64) bool {
	restaurantsClosingTime := helpers.ConvertStringTimeToDesiredUnixFormat(response.Openingtime.Restauranttime.Ranges[0].End)
	kitchenClosingTime := helpers.ConvertStringTimeToDesiredUnixFormat(response.Openingtime.Kitchentime.Ranges[0].End)

	if currentTime > restaurantsClosingTime || currentTime > kitchenClosingTime {
		return true
	}

	return false
}

func (response *ResponseFields) reservationLinkIsNotValid() bool {
	return !strings.Contains(response.Links.TableReservationLocalized.FiFi, "https://s-varaukset.fi/online/reservation/fi/")
}

func (response *ResponseFields) cityDoesNotMatchUsersCity(usersCity string) bool {
	response.Address.Municipality.FiFi = strings.ToLower(response.Address.Municipality.FiFi)
	restaurantsCity := response.Address.Municipality.FiFi

	return restaurantsCity != strings.ToLower(usersCity)
}
