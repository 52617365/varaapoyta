package raflaamoRestaurantsApi

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"strings"
)

type RaflaamoRestaurantsApi struct {
	httpClient *http.Client
	request    *http.Request
	response   *http.Response
	usersCity  string
}

func getRaflaamoRestaurantsApi(city string) (*RaflaamoRestaurantsApi, error) {
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
		usersCity:  city,
	}, nil
}

func (restaurantsApi *RaflaamoRestaurantsApi) getRestaurantsFromRaflaamoApi(city string) ([]responseFields, error) {
	restaurantsApi, err := getRaflaamoRestaurantsApi(city)
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

	validRestaurantsMatchingCriteria := restaurantsApi.filterBadRestaurantsOut(decodedRaflaamoRestaurants)

	// TODO: Handle return value being empty in caller.
	return validRestaurantsMatchingCriteria, nil
}

//
// 	A restaurant is considered "Bad" if:
//		- Restaurants city is not from the provided city.
//		- Restaurants reservation link does not exist or contains odd contents.
//		- Restaurant does not contain opening times (Specified in the Ranges array).
//
func (restaurantsApi *RaflaamoRestaurantsApi) filterBadRestaurantsOut(structureContainingRestaurantData *responseTopLevel) []responseFields {
	restaurantsApi.usersCity = strings.ToLower(restaurantsApi.usersCity)
	arrayContainingRestaurantData := structureContainingRestaurantData.Data.ListRestaurantsByLocation.Edges

	filteredRestaurantsFromProvidedCity := make([]responseFields, 0, 50)
	for _, restaurant := range arrayContainingRestaurantData {
		if restaurant.isBad(restaurantsApi.usersCity) {
			continue
		}

		// Here we have done all of the checks we know to date. There might be more in the future once I figure them out.
		filteredRestaurantsFromProvidedCity = append(filteredRestaurantsFromProvidedCity, restaurant)
	}
	return filteredRestaurantsFromProvidedCity
}

func (restaurant *responseFields) isBad(city string) bool {
	if restaurant.cityDoesNotMatchUsersCity(city) {
		return true
	}
	if restaurant.reservationLinkIsNotValid() {
		return true
	}
	if restaurant.doesNotContainOpeningTimes() {
		return true
	}
	return false
}

func (response *responseFields) doesNotContainOpeningTimes() bool {
	restaurantsOpeningTimes := response.Openingtime.Restauranttime.Ranges
	kitchensOpeningTimes := response.Openingtime.Kitchentime.Ranges

	if restaurantsOpeningTimes == nil || kitchensOpeningTimes == nil {
		return true
	}
	return false
}

func (response *responseFields) reservationLinkIsNotValid() bool {
	return !strings.Contains(response.Links.TableReservationLocalized.FiFi, "https://s-varaukset.fi/online/reservation/fi/")
}

func (response *responseFields) cityDoesNotMatchUsersCity(usersCity string) bool {
	response.Name.FiFi = strings.ToLower(response.Name.FiFi)
	restaurantsName := response.Name.FiFi

	return restaurantsName != usersCity
}
