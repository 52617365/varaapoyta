package main

import (
	raflaamoGraphApi3 "backend/raflaamoGraphApi"
	raflaamoRestaurantsApi2 "backend/raflaamoRestaurantsApi"
	"backend/timeUtils"
	"regexp"
)

// TODO: iterate restaurants and call getGraphApiResponseFromTimeSlot on each one.
// TODO: regexes should be compiled in main.go

type ResponseFields = raflaamoRestaurantsApi2.ResponseFields
type RaflaamoRestaurantsApi = raflaamoRestaurantsApi2.RaflaamoRestaurantsApi

// All the regex we want compiled in main.go.
// regexToMatchRestaurantId := regexp.MustCompile(`[^fi/]\d+`)
// regexToMatchTime := regexp.MustCompile(`\d{2}:\d{2}`)
// regexToMatchDate := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
func getAvailableTablesFromRestaurants(regexToMatchRestaurantId regexp.Regexp, regexToMatchTime *regexp.Regexp, regexToMatchDate *regexp.Regexp, city string, amountOfEaters int) error {
	raflaamoRelatedTimes := timeUtils.GetRaflaamoTimes(regexToMatchTime, regexToMatchDate)
	raflaamoGraphApi := raflaamoGraphApi3.GetRaflaamoGraphApi()
	raflaamoRestaurantsApi, err := raflaamoRestaurantsApi2.GetRaflaamoRestaurantsApiStruct(city)
	if err != nil {
		return err
	}
	restaurants, err := raflaamoRestaurantsApi.GetRestaurants()
	if err != nil {
		return err
	}
	// TODO: Do something with restaurants.
	//for _, restaurant := range restaurants {
	//	raflaamoGraphApiPayload := raflaamoGraphApi.GetRaflaamoGraphApiPayload(restaurant.Links.TableReservationLocalized.FiFi, amountOfEaters, raflaamoRelatedTimes.TimeAndDate.CurrentDate, &regexToMatchRestaurantId)
	//	restaurantGraphApiRequestUrls := raflaamoGraphApiPayload.IterateAllPossibleTimeSlotsAndGenerateRequestUrls(raflaamoRelatedTimes)
	//
	//}

	//raflaamoTimes.getAllGraphApiUnixTimeIntervalsFromCurrentPointForward(restaurantClosingTimeUnix) // This should be called in caller for each restaurant because closing times change.
}
