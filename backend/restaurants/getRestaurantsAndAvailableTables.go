package restaurants

import (
	"backend/raflaamoGraphApi"
	"backend/raflaamoRestaurantsApi"
	"backend/raflaamoTime"
	"fmt"
)

// GetRestaurantsAndAvailableTables This is the entry point to the functionality.
func GetRestaurantsAndAvailableTables(city string, amountOfEaters int) error {
	allNeededRaflaamoTimes := raflaamoTime.GetAllNeededRaflaamoTimes(regexToMatchTime, regexToMatchDate)
	graphApi := raflaamoGraphApi.GetRaflaamoGraphApi()
	restaurantsApi, err := raflaamoRestaurantsApi.GetRaflaamoRestaurantsApi(city)
	if err != nil {
		return err
	}

	allRestaurantsFromRaflaamoRestaurantsApi, err := restaurantsApi.GetAllRestaurantsFromRaflaamoRestaurantsApi()
	if err != nil {
		return err
	}

	// GetRestaurantsAndAvailableTables TODO: Use worker-pool here to speed stuff up.
	for _, restaurant := range allRestaurantsFromRaflaamoRestaurantsApi {
		openTablesForRestaurant, _ := getAvailableTablesForRestaurant(&restaurant, allNeededRaflaamoTimes, amountOfEaters, graphApi)
		// GetRestaurantsAndAvailableTables TODO: store openTablesForRestaurant somewhere.
		fmt.Println(openTablesForRestaurant)
	}
	return nil
}

func getAvailableTablesForRestaurant(restaurant *raflaamoRestaurantsApi.ResponseFields, raflaamoRelatedTimes *raflaamoTime.RaflaamoTimes, amountOfEaters int, graphApi *RaflaamoGraphApi) ([]*parsedGraphData, error) {
	restaurantsKitchenClosingTime := restaurant.Openingtime.Kitchentime.Ranges[0].End
	raflaamoRelatedTimes.GetAllGraphApiUnixTimeIntervalsFromCurrentPointForward(restaurantsKitchenClosingTime)

	raflaamoGraphApiRequestUrlStruct := raflaamoGraphApi.GetRaflaamoGraphApiRequestUrl(restaurant.Links.TableReservationLocalized.FiFi, amountOfEaters, raflaamoRelatedTimes.TimeAndDate.CurrentDate, regexToMatchRestaurantId)

	restaurantGraphApiRequestUrls := raflaamoGraphApiRequestUrlStruct.GenerateGraphApiRequestUrlsForRestaurant(raflaamoRelatedTimes)

	availableTablesFromRestaurant, err := getAvailableTableTimesFromRestaurantRequestUrls(restaurant, raflaamoRelatedTimes, restaurantGraphApiRequestUrls, graphApi)

	//restaurant.AvailableTimeSlots = availableTablesFromRestaurant
	if err != nil {
		return nil, err
	}

	return availableTablesFromRestaurant, nil
}

func getAvailableTableTimesFromRestaurantRequestUrls(restaurant *raflaamoRestaurantsApi.ResponseFields, raflaamoRelatedTimes *raflaamoTime.RaflaamoTimes, restaurantGraphApiRequestUrls []string, graphApi *RaflaamoGraphApi) ([]*parsedGraphData, error) {
	parsedGraphApiResponses := make([]*parsedGraphData, 0, len(restaurantGraphApiRequestUrls))

	for _, restaurantGraphApiRequestUrl := range restaurantGraphApiRequestUrls {
		graphApiResponseFromRequestUrl, err := graphApi.GetGraphApiResponseFromTimeSlot(restaurantGraphApiRequestUrl)
		if err != nil {
			return nil, err
		}

		// getAvailableTableTimesFromRestaurantRequestUrls TODO: capture all the available time slots from the response intervals. 11.15 11.30 etc.
		// getAvailableTableTimesFromRestaurantRequestUrls TODO: remember to take into consideration the kitchens closing time (can't reserve 1h before kitchen closes.)
		// getAvailableTableTimesFromRestaurantRequestUrls TODO: capture restaurants time till kitchen and restaurant closes.
		parsedGraphApiResponses = append(parsedGraphApiResponses, graphApiResponseFromRequestUrl)
	}
	return parsedGraphApiResponses, nil
}
