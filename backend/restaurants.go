package main

import (
	"backend/raflaamoRestaurantsApi"
	"fmt"
	"regexp"
)

// TODO: iterate restaurants and call getGraphApiResponseFromTimeSlot on each one.

// All the regex we want compiled in main.go.
// regexToMatchRestaurantId := regexp.MustCompile(`[^fi/]\d+`)
// regexToMatchTime := regexp.MustCompile(`\d{2}:\d{2}`)
// regexToMatchDate := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
func getAvailableTablesFromRestaurants(regexToMatchRestaurantId *regexp.Regexp, regexToMatchTime *regexp.Regexp, regexToMatchDate *regexp.Regexp, city string, amountOfEaters int) error {
	//raflaamoRelatedTimes := timeUtils.GetRaflaamoTimes(regexToMatchTime, regexToMatchDate)
	//raflaamoGraphApi := raflaamoGraphApi.GetRaflaamoGraphApi()
	restaurantsApi, err := raflaamoRestaurantsApi.GetRaflaamoRestaurantsApiStruct(city)
	if err != nil {
		return err
	}
	restaurants, err := restaurantsApi.GetRestaurants()
	if err != nil {
		return err
	}
	for _, restaurant := range restaurants {
		fmt.Println(restaurant.Name.FiFi)
	}
	// TODO: Do something with restaurants.
	//for _, restaurant := range restaurants {
	//	raflaamoGraphApiPayload := raflaamoGraphApi.GetRaflaamoGraphApiPayload(restaurant.Links.TableReservationLocalized.FiFi, amountOfEaters, raflaamoRelatedTimes.TimeAndDate.CurrentDate, regexToMatchRestaurantId)
	//	restaurantGraphApiRequestUrls := raflaamoGraphApiPayload.IterateAllPossibleTimeSlotsAndGenerateRequestUrls(raflaamoRelatedTimes)
	//
	//}

	//raflaamoTimes.getAllGraphApiUnixTimeIntervalsFromCurrentPointForward(restaurantClosingTimeUnix) // This should be called in caller for each restaurant because closing times change.
	return nil
}
