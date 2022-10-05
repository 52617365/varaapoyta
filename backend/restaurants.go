package main

import (
	"backend/raflaamoGraphApi"
	"backend/raflaamoRestaurantsApi"
	"backend/timeUtils"
	"fmt"
	"regexp"
)

// TODO: iterate restaurants and call getGraphApiResponseFromTimeSlot on each one.

var regexToMatchRestaurantId = regexp.MustCompile(`[^fi/]\d+`)
var regexToMatchTime = regexp.MustCompile(`\d{2}:\d{2}`)
var regexToMatchDate = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)

func getAvailableTablesFromRestaurants(city string, amountOfEaters int) error {
	raflaamoRelatedTimes := timeUtils.GetRaflaamoTimes(regexToMatchTime, regexToMatchDate)
	//graphApi := raflaamoGraphApi.GetRaflaamoGraphApi()
	restaurantsApi, err := raflaamoRestaurantsApi.GetRaflaamoRestaurantsApiStruct(city)
	if err != nil {
		return err
	}
	restaurants, err := restaurantsApi.GetRestaurants()
	if err != nil {
		return err
	}

	for _, restaurant := range restaurants {
		restaurantClosingTime := restaurant.Openingtime.Kitchentime.Ranges[0].End
		raflaamoRelatedTimes.GetAllGraphApiUnixTimeIntervalsFromCurrentPointForward(restaurantClosingTime)
		raflaamoGraphApiPayload := raflaamoGraphApi.GetRaflaamoGraphApiPayload(restaurant.Links.TableReservationLocalized.FiFi, amountOfEaters, raflaamoRelatedTimes.TimeAndDate.CurrentDate, regexToMatchRestaurantId)
		restaurantGraphApiRequestUrls := raflaamoGraphApiPayload.IterateAllPossibleTimeSlotsAndGenerateRequestUrls(raflaamoRelatedTimes)

		for _, restaurantGraphApiRequestUrl := range restaurantGraphApiRequestUrls {
			fmt.Println(restaurantGraphApiRequestUrl)
			// TODO: send graph api requests here.
		}

	}

	//raflaamoTimes.getAllGraphApiUnixTimeIntervalsFromCurrentPointForward(restaurantClosingTimeUnix) // This should be called in caller for each restaurant because closing times change.
	return nil
}
