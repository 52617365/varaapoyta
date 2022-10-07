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
		getAvailableTablesForRestaurant(restaurant, allNeededRaflaamoTimes, amountOfEaters, initializedRaflaamoGraphApi)
	}
	return allRestaurantsFromRaflaamoRestaurantsApi, nil
}

func getAvailableTablesForRestaurant(restaurant *raflaamoRestaurantsApi.ResponseFields, raflaamoRelatedTimes *raflaamoTime.RaflaamoTimes, amountOfEaters int, graphApi *raflaamoGraphApi.RaflaamoGraphApi) {
	restaurantsKitchenClosingTime := restaurant.Openingtime.Kitchentime.Ranges[0].End
	raflaamoRelatedTimes.GetAllGraphApiUnixTimeIntervalsFromCurrentPointForward(restaurantsKitchenClosingTime)

	raflaamoGraphApiRequestUrlStruct := raflaamoGraphApi.GetRaflaamoGraphApiRequestUrl(restaurant.Links.TableReservationLocalized.FiFi, amountOfEaters, raflaamoRelatedTimes.TimeAndDate.CurrentDate, regexToMatchRestaurantId)
	restaurant.Links.TableReservationLocalizedId = raflaamoGraphApiRequestUrlStruct.IdFromReservationPageUrl // Storing the id for the front end.

	restaurantGraphApiRequestUrls := raflaamoGraphApiRequestUrlStruct.GenerateGraphApiRequestUrlsForRestaurant(raflaamoRelatedTimes)

	getAvailableTableTimesFromRestaurantRequestUrls(restaurant, raflaamoRelatedTimes, restaurantGraphApiRequestUrls, graphApi)
}

func getAvailableTableTimesFromRestaurantRequestUrls(restaurant *raflaamoRestaurantsApi.ResponseFields, raflaamoRelatedTimes *raflaamoTime.RaflaamoTimes, restaurantGraphApiRequestUrls []string, graphApi *raflaamoGraphApi.RaflaamoGraphApi) {
	//restaurant.GraphApiResults.AvailableTimeSlotsBuffer = make(chan string, 96) // 96 is worst case scenario.
	//restaurant.AvailableTimeSlots = make([]string, 0, len(raflaamoRelatedTimes.AllRaflaamoReservationTimeIntervals))

	graphApiResultsFromRequestUrls := raflaamoRestaurantsApi.GraphApiResult{AvailableTimeSlotsBuffer: make(chan string, 96), Err: make(chan error, 96)}
	for _, restaurantGraphApiRequestUrl := range restaurantGraphApiRequestUrls {
		restaurantGraphApiRequestUrl := restaurantGraphApiRequestUrl
		go func() {
			graphApiResponseFromRequestUrl, err := graphApi.GetGraphApiResponseFromTimeSlot(restaurantGraphApiRequestUrl)
			if err != nil {
				restaurant.GraphApiResults.Err <- err
				return
			}

			// getAvailableTableTimesFromRestaurantRequestUrls TODO: remember to take into consideration the kitchens closing time (can't reserve 1h before kitchen closes.)
			// getAvailableTableTimesFromRestaurantRequestUrls TODO: capture restaurants time till kitchen and restaurant closes.
			graphApiReservationTimes := raflaamoTime.GetGraphApiReservationTimes(graphApiResponseFromRequestUrl)

			timeSlotIntervalsFromRestaurant := graphApiReservationTimes.GetTimeSlotsInBetweenIntervals(raflaamoRelatedTimes.AllRaflaamoReservationTimeIntervals)
			for _, timeSlotInterval := range timeSlotIntervalsFromRestaurant {
				graphApiResultsFromRequestUrls.AvailableTimeSlotsBuffer <- timeSlotInterval
			}
		}()
	}
	close(graphApiResultsFromRequestUrls.AvailableTimeSlotsBuffer)
	close(graphApiResultsFromRequestUrls.Err)
	restaurant.GraphApiResults = &graphApiResultsFromRequestUrls
}
