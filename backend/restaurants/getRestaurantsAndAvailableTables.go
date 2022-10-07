package restaurants

import (
	"backend/raflaamoGraphApi"
	"backend/raflaamoRestaurantsApi"
	"backend/raflaamoTime"
)

// GetRestaurantsAndAvailableTables This is the entry point to the functionality.
func GetRestaurantsAndAvailableTables(city string, amountOfEaters int) ([]raflaamoRestaurantsApi.ResponseFields, error) {
	allNeededRaflaamoTimes := raflaamoTime.GetAllNeededRaflaamoTimes(regexToMatchTime, regexToMatchDate)
	initializedRaflaamoGraphApi := raflaamoGraphApi.GetRaflaamoGraphApi()
	initializedRaflaamoRestaurantsApi, err := raflaamoRestaurantsApi.GetRaflaamoRestaurantsApi(city)
	if err != nil {
		return nil, err
	}

	allRestaurantsFromRaflaamoRestaurantsApi, err := initializedRaflaamoRestaurantsApi.GetAllRestaurantsFromRaflaamoRestaurantsApi()
	if err != nil {
		return nil, err
	}

	// TODO: use goroutines here to speed stuff up.
	for index, _ := range allRestaurantsFromRaflaamoRestaurantsApi {
		restaurant := &allRestaurantsFromRaflaamoRestaurantsApi[index]
		err := getAvailableTablesForRestaurant(restaurant, allNeededRaflaamoTimes, amountOfEaters, initializedRaflaamoGraphApi)
		if err != nil {
			continue
		}
	}
	return allRestaurantsFromRaflaamoRestaurantsApi, nil
}

func getAvailableTablesForRestaurant(restaurant *raflaamoRestaurantsApi.ResponseFields, raflaamoRelatedTimes *raflaamoTime.RaflaamoTimes, amountOfEaters int, graphApi *raflaamoGraphApi.RaflaamoGraphApi) error {
	restaurantsKitchenClosingTime := restaurant.Openingtime.Kitchentime.Ranges[0].End
	raflaamoRelatedTimes.GetAllGraphApiUnixTimeIntervalsFromCurrentPointForward(restaurantsKitchenClosingTime)

	raflaamoGraphApiRequestUrlStruct := raflaamoGraphApi.GetRaflaamoGraphApiRequestUrl(restaurant.Links.TableReservationLocalized.FiFi, amountOfEaters, raflaamoRelatedTimes.TimeAndDate.CurrentDate, regexToMatchRestaurantId)
	restaurant.Links.TableReservationLocalizedId = raflaamoGraphApiRequestUrlStruct.IdFromReservationPageUrl // Storing the id for the front end.

	restaurantGraphApiRequestUrls := raflaamoGraphApiRequestUrlStruct.GenerateGraphApiRequestUrlsForRestaurant(raflaamoRelatedTimes)

	err := getAvailableTableTimesFromRestaurantRequestUrls(restaurant, raflaamoRelatedTimes, restaurantGraphApiRequestUrls, graphApi)

	if err != nil {
		return err
	}

	return nil
}

func getAvailableTableTimesFromRestaurantRequestUrls(restaurant *raflaamoRestaurantsApi.ResponseFields, raflaamoRelatedTimes *raflaamoTime.RaflaamoTimes, restaurantGraphApiRequestUrls []string, graphApi *raflaamoGraphApi.RaflaamoGraphApi) error {
	restaurant.AvailableTimeSlots = make([]string, 0, len(raflaamoRelatedTimes.AllRaflaamoReservationTimeIntervals))

	for _, restaurantGraphApiRequestUrl := range restaurantGraphApiRequestUrls {
		graphApiResponseFromRequestUrl, err := graphApi.GetGraphApiResponseFromTimeSlot(restaurantGraphApiRequestUrl)
		if err != nil {
			return err
		}

		// getAvailableTableTimesFromRestaurantRequestUrls TODO: remember to take into consideration the kitchens closing time (can't reserve 1h before kitchen closes.)
		// getAvailableTableTimesFromRestaurantRequestUrls TODO: capture restaurants time till kitchen and restaurant closes.
		graphApiReservationTimes := raflaamoTime.GetGraphApiReservationTimes(graphApiResponseFromRequestUrl)

		timeSlotIntervalsFromRestaurant := graphApiReservationTimes.GetTimeSlotsInBetweenIntervals(raflaamoRelatedTimes.AllRaflaamoReservationTimeIntervals)
		restaurant.AvailableTimeSlots = append(restaurant.AvailableTimeSlots, timeSlotIntervalsFromRestaurant...) // @notice will be channelled in the future
	}
	return nil
}
