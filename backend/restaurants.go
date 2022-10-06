package main

import (
	"backend/raflaamoGraphApi"
	"backend/raflaamoRestaurantsApi"
	"backend/timeUtils"
	"regexp"
)

var regexToMatchRestaurantId = regexp.MustCompile(`[^fi/]\d+`)
var regexToMatchTime = regexp.MustCompile(`\d{2}:\d{2}`)
var regexToMatchDate = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)

// GetRestaurantsAndAvailableTables TODO: Use goroutines to speed stuff up.
func GetRestaurantsAndAvailableTables(city string, amountOfEaters int) error {
	allNeededRaflaamoTimes := timeUtils.GetAllNeededRaflaamoTimes(regexToMatchTime, regexToMatchDate)
	graphApi := raflaamoGraphApi.GetRaflaamoGraphApi()
	restaurantsApi, err := raflaamoRestaurantsApi.GetRaflaamoRestaurantsApi(city)
	if err != nil {
		return err
	}

	allRestaurantsFromRaflaamoRestaurantsApi, err := restaurantsApi.GetAllRestaurantsFromRaflaamoRestaurantsApi()
	if err != nil {
		return err
	}

	for _, restaurant := range allRestaurantsFromRaflaamoRestaurantsApi {
		openTablesForRestaurant := getAvailableTablesForRestaurant(&restaurant, allNeededRaflaamoTimes, amountOfEaters, graphApi)
		// GetRestaurantsAndAvailableTables TODO: store openTablesForRestaurant somewhere.

	}
	return nil
}

type RaflaamoGraphApi = raflaamoGraphApi.RaflaamoGraphApi

func getAvailableTablesForRestaurant(restaurant *raflaamoRestaurantsApi.ResponseFields, raflaamoRelatedTimes *timeUtils.RaflaamoTimes, amountOfEaters int, graphApi *RaflaamoGraphApi) ([]*parsedGraphData, error) {
	restaurantsKitchenClosingTime := restaurant.Openingtime.Kitchentime.Ranges[0].End
	raflaamoRelatedTimes.GetAllGraphApiUnixTimeIntervalsFromCurrentPointForward(restaurantsKitchenClosingTime)

	raflaamoGraphApiRequestUrlStruct := raflaamoGraphApi.GetRaflaamoGraphApiRequestUrl(restaurant.Links.TableReservationLocalized.FiFi, amountOfEaters, raflaamoRelatedTimes.TimeAndDate.CurrentDate, regexToMatchRestaurantId)

	restaurantGraphApiRequestUrls := raflaamoGraphApiRequestUrlStruct.GenerateGraphApiRequestUrlsForRestaurant(raflaamoRelatedTimes)

	availableTablesFromRestaurant, err := getAvailableTablesFromRestaurantRequestUrls(restaurantGraphApiRequestUrls, graphApi)

	if err != nil {
		return nil, err
	}

	return availableTablesFromRestaurant, nil
}

type parsedGraphData = raflaamoGraphApi.ParsedGraphData

func getAvailableTablesFromRestaurantRequestUrls(restaurantGraphApiRequestUrls []string, graphApi *RaflaamoGraphApi) ([]*parsedGraphData, error) {
	parsedGraphApiResponses := make([]*parsedGraphData, 0, len(restaurantGraphApiRequestUrls))

	for _, restaurantGraphApiRequestUrl := range restaurantGraphApiRequestUrls {
		graphApiResponseFromRequestUrl, err := graphApi.GetGraphApiResponseFromTimeSlot(restaurantGraphApiRequestUrl)
		if err != nil {
			return nil, err
		}

		// GetRestaurantsAndAvailableTables TODO: capture all the available time slots from the response intervals. 11.15 11.30 etc.
		// GetRestaurantsAndAvailableTables TODO: remember to take into consideration the kitchens closing time (can't reserve 1h before kitchen closes.)
		// GetRestaurantsAndAvailableTables TODO: capture restaurants time till kitchen and restaurant closes.
		parsedGraphApiResponses = append(parsedGraphApiResponses, graphApiResponseFromRequestUrl)
	}
	return parsedGraphApiResponses, nil
}
