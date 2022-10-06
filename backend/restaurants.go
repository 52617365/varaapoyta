package main

import (
	"backend/raflaamoGraphApi"
	"backend/raflaamoRestaurantsApi"
	"backend/timeUtils"
	"fmt"
	"regexp"
)

var regexToMatchRestaurantId = regexp.MustCompile(`[^fi/]\d+`)
var regexToMatchTime = regexp.MustCompile(`\d{2}:\d{2}`)
var regexToMatchDate = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)

func IterateAllRestaurants(city string, amountOfEaters int) error {
	raflaamoRelatedTimes := timeUtils.GetRaflaamoTimes(regexToMatchTime, regexToMatchDate)
	restaurantsApi, err := raflaamoRestaurantsApi.GetRaflaamoRestaurantsApiStruct(city)
	graphApi := raflaamoGraphApi.GetRaflaamoGraphApi()
	if err != nil {
		return err
	}
	restaurants, err := restaurantsApi.GetRestaurants()
	if err != nil {
		return err
	}

	for _, restaurant := range restaurants {
		restaurantOpenTables := getAvailableTablesForRestaurant(&restaurant, raflaamoRelatedTimes, amountOfEaters, graphApi)
		// TODO: store it somewhere

	}

	//raflaamoTimes.getAllGraphApiUnixTimeIntervalsFromCurrentPointForward(restaurantClosingTimeUnix) // This should be called in caller for each restaurant because closing times change.
	return nil
}

type RaflaamoGraphApi = raflaamoGraphApi.RaflaamoGraphApi

func getAvailableTablesForRestaurant(restaurant *raflaamoRestaurantsApi.ResponseFields, raflaamoRelatedTimes *timeUtils.RaflaamoTimes, amountOfEaters int, graphApi *RaflaamoGraphApi) []*parsedGraphData {
	restaurantClosingTime := restaurant.Openingtime.Kitchentime.Ranges[0].End
	raflaamoRelatedTimes.GetAllGraphApiUnixTimeIntervalsFromCurrentPointForward(restaurantClosingTime)
	raflaamoGraphApiPayload := raflaamoGraphApi.GetRaflaamoGraphApiPayload(restaurant.Links.TableReservationLocalized.FiFi, amountOfEaters, raflaamoRelatedTimes.TimeAndDate.CurrentDate, regexToMatchRestaurantId)
	restaurantGraphApiRequestUrls := raflaamoGraphApiPayload.IterateAllPossibleTimeSlotsAndGenerateRequestUrls(raflaamoRelatedTimes)

	availableTablesFromRestaurant := getAvailableTablesFromRestaurantRequestUrls(restaurantGraphApiRequestUrls, graphApi)
	return availableTablesFromRestaurant
}

type parsedGraphData = raflaamoGraphApi.ParsedGraphData

func getAvailableTablesFromRestaurantRequestUrls(restaurantGraphApiRequestUrls []string, graphApi *RaflaamoGraphApi) []*parsedGraphData {
	parsedGraphApiResponses := make([]*parsedGraphData, 0, len(restaurantGraphApiRequestUrls))

	for _, restaurantGraphApiRequestUrl := range restaurantGraphApiRequestUrls {
		fmt.Println("restaurant api request url is", restaurantGraphApiRequestUrl)
		graphApiResponseFromRequestUrl, err := graphApi.GetGraphApiResponseFromTimeSlot(restaurantGraphApiRequestUrl)
		if err != nil {
			return nil
		}
		parsedGraphApiResponses = append(parsedGraphApiResponses, graphApiResponseFromRequestUrl)
	}
	return parsedGraphApiResponses
}
