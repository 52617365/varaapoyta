package restaurants

import (
	"backend/raflaamoGraphApi"
	"backend/raflaamoRestaurantsApi"
	"backend/raflaamoTime"
	"fmt"
	"golang.org/x/exp/slices"
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

func getAvailableTablesForRestaurant(restaurant *raflaamoRestaurantsApi.ResponseFields, raflaamoRelatedTimes *raflaamoTime.RaflaamoTimes, amountOfEaters int, graphApi *raflaamoGraphApi.RaflaamoGraphApi) ([]string, error) {
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

func getAvailableTableTimesFromRestaurantRequestUrls(restaurant *raflaamoRestaurantsApi.ResponseFields, raflaamoRelatedTimes *raflaamoTime.RaflaamoTimes, restaurantGraphApiRequestUrls []string, graphApi *raflaamoGraphApi.RaflaamoGraphApi) ([]string, error) {
	parsedGraphApiTimeSlots := make([]string, 0, len(restaurantGraphApiRequestUrls))

	for _, restaurantGraphApiRequestUrl := range restaurantGraphApiRequestUrls {
		graphApiResponseFromRequestUrl, err := graphApi.GetGraphApiResponseFromTimeSlot(restaurantGraphApiRequestUrl)
		if err != nil {
			return nil, err
		}

		// getAvailableTableTimesFromRestaurantRequestUrls TODO: remember to take into consideration the kitchens closing time (can't reserve 1h before kitchen closes.)
		// getAvailableTableTimesFromRestaurantRequestUrls TODO: capture restaurants time till kitchen and restaurant closes.
		graphApiReservationTimes := raflaamoTime.GetGraphApiReservationTimes(graphApiResponseFromRequestUrl)

		timeSlotIntervalsFromRestaurant := graphApiReservationTimes.GetTimeSlotsInBetweenIntervals(raflaamoRelatedTimes.AllRaflaamoReservationTimeIntervals)
		restaurant.AvailableTimeSlots = timeSlotIntervalsFromRestaurant
		for _, timeSlotInterval := range timeSlotIntervalsFromRestaurant {
			if !slices.Contains(parsedGraphApiTimeSlots, timeSlotInterval) {
				parsedGraphApiTimeSlots = append(parsedGraphApiTimeSlots, timeSlotInterval)
			}
		}
	}
	return parsedGraphApiTimeSlots, nil
}
