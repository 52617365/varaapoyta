package raflaamoRestaurantsApi

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func GetRaflaamoRestaurantsApiStruct(city string) (*RaflaamoRestaurantsApi, error) {
	httpClient := &http.Client{}
	const data = `{"operationName":"getRestaurantsByLocation","variables":{"first":1000,"input":{"restaurantType":"ALL","locationName":"Helsinki","feature":{"rentableVenues":false}},"after":"eyJmIjowLCJnIjp7ImEiOjYwLjE3MTE2LCJvIjoyNC45MzI1OH19"},"query":"fragment Locales on LocalizedString {fi_FI }fragment Restaurant on Restaurant {  id  name {    ...Locales    }  address {    municipality {      ...Locales       }        street {      ...Locales       }       zipCode     }   openingTime {    restaurantTime {      ranges {        start        end             }             }    kitchenTime {      ranges {        start        end        endNextDay              }             }    }  links {    tableReservationLocalized {      ...Locales        }    homepageLocalized {      ...Locales          }   }     }query getRestaurantsByLocation($first: Int, $after: String, $input: ListRestaurantsByLocationInput!) {  listRestaurantsByLocation(first: $first, after: $after, input: $input) {    totalCount      edges {      ...Restaurant        }     }}"}`

	req, err := http.NewRequest("POST", "https://api.raflaamo.fi/query", bytes.NewBuffer([]byte(data)))

	if err != nil {
		return nil, fmt.Errorf("[GetRaflaamoRestaurantsApiStruct] - %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("client_id", "jNAWMvWD9rp637RaR")
	req.Header.Add("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	return &RaflaamoRestaurantsApi{
		httpClient:               httpClient,
		request:                  req,
		cityToGetRestaurantsFrom: city,
	}, nil
}

func (raflaamoRestaurantsApi *RaflaamoRestaurantsApi) getRestaurantsFromRaflaamoApi() ([]ResponseFields, error) {
	httpClient := raflaamoRestaurantsApi.httpClient
	request := raflaamoRestaurantsApi.request

	res, err := httpClient.Do(request)

	if err != nil {
		return nil, fmt.Errorf("[GetRestaurantsFromRaflaamoApi] - %w", errors.New("there was an error connecting to the raflaamo api"))
	}

	raflaamoRestaurantsApi.response = res

	decodedRaflaamoRestaurants, err := raflaamoRestaurantsApi.deserializeRaflaamoRestaurantsResponse()

	if err != nil {
		return nil, fmt.Errorf("[GetRestaurantsFromRaflaamoApi] - %w", errors.New("there was an error deserializing raflaamo API response"))
	}

	validRestaurantsMatchingCriteria := raflaamoRestaurantsApi.filterBadRestaurantsOut(decodedRaflaamoRestaurants)

	// TODO: Handle return value being empty in caller.
	return validRestaurantsMatchingCriteria, nil
}

// A restaurant is considered "Bad" if:
//   - Restaurants city is not from the provided city.
//   - Restaurants reservation link does not exist or contains odd contents.
//   - Restaurant does not contain opening times (Specified in the Ranges array).
func (raflaamoRestaurantsApi *RaflaamoRestaurantsApi) filterBadRestaurantsOut(structureContainingRestaurantData *responseTopLevel) []ResponseFields {
	raflaamoRestaurantsApi.cityToGetRestaurantsFrom = strings.ToLower(raflaamoRestaurantsApi.cityToGetRestaurantsFrom)
	arrayContainingRestaurantData := structureContainingRestaurantData.Data.ListRestaurantsByLocation.Edges

	filteredRestaurantsFromProvidedCity := make([]ResponseFields, 0, 50)
	for _, restaurant := range arrayContainingRestaurantData {
		if restaurant.isBad(raflaamoRestaurantsApi.cityToGetRestaurantsFrom) {
			continue
		}

		// Here we have done all the checks we know to date. There might be more in the future once I figure them out.
		filteredRestaurantsFromProvidedCity = append(filteredRestaurantsFromProvidedCity, restaurant)
	}
	return filteredRestaurantsFromProvidedCity
}

func (response *ResponseFields) isBad(city string) bool {
	if response.cityDoesNotMatchUsersCity(city) {
		return true
	}
	if response.reservationLinkIsNotValid() {
		return true
	}
	if response.doesNotContainOpeningTimes() {
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

func (response *ResponseFields) reservationLinkIsNotValid() bool {
	return !strings.Contains(response.Links.TableReservationLocalized.FiFi, "https://s-varaukset.fi/online/reservation/fi/")
}

func (response *ResponseFields) cityDoesNotMatchUsersCity(usersCity string) bool {
	response.Address.Municipality.FiFi = strings.ToLower(response.Address.Municipality.FiFi)
	restaurantsCity := response.Address.Municipality.FiFi

	return restaurantsCity != usersCity
}

func (raflaamoRestaurantsApi *RaflaamoRestaurantsApi) GetRestaurants() ([]ResponseFields, error) {
	restaurantsApi, err := GetRaflaamoRestaurantsApiStruct(raflaamoRestaurantsApi.cityToGetRestaurantsFrom)
	if err != nil {
		return nil, err
	}

	restaurantsFromApi, err := restaurantsApi.getRestaurantsFromRaflaamoApi()
	if err != nil {
		return nil, err
	}

	return restaurantsFromApi, nil
}
