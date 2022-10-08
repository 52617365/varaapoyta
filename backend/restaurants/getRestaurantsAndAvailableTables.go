package restaurants

import (
	"backend/raflaamoGraphApi"
	"backend/raflaamoRestaurantsApi"
	"backend/raflaamoTime"
	"fmt"
	"sync"
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
		go getAvailableTablesForRestaurant(restaurant, allNeededRaflaamoTimes, amountOfEaters, initializedRaflaamoGraphApi)
	}
	return allRestaurantsFromRaflaamoRestaurantsApi, nil
}

func getAvailableTablesForRestaurant(restaurant *raflaamoRestaurantsApi.ResponseFields, raflaamoRelatedTimes *raflaamoTime.RaflaamoTimes, amountOfEaters int, graphApi *raflaamoGraphApi.RaflaamoGraphApi) {
	restaurantsKitchenClosingTime := restaurant.Openingtime.Kitchentime.Ranges[0].End
	raflaamoRelatedTimes.GetAllGraphApiUnixTimeIntervalsFromCurrentPointForward(restaurantsKitchenClosingTime)

	raflaamoGraphApiRequestUrlStruct := raflaamoGraphApi.GetRaflaamoGraphApiRequestUrl(restaurant.Links.TableReservationLocalized.FiFi, amountOfEaters, raflaamoRelatedTimes.TimeAndDate.CurrentDate, regexToMatchRestaurantId)
	restaurant.Links.TableReservationLocalizedId = raflaamoGraphApiRequestUrlStruct.IdFromReservationPageUrl // Storing the id for the front end.

	restaurantGraphApiRequestUrls := raflaamoGraphApiRequestUrlStruct.GenerateGraphApiRequestUrlsForRestaurant(raflaamoRelatedTimes)

	getAvailableTableTimesFromRestaurantRequestUrlsIntoRestaurantsChannel(restaurant, raflaamoRelatedTimes, restaurantGraphApiRequestUrls, graphApi)
}

// TODO: make the channel stuff work here.
func getAvailableTableTimesFromRestaurantRequestUrlsIntoRestaurantsChannel(restaurant *raflaamoRestaurantsApi.ResponseFields, raflaamoRelatedTimes *raflaamoTime.RaflaamoTimes, restaurantGraphApiRequestUrls []string, graphApi *raflaamoGraphApi.RaflaamoGraphApi) {
	var wg sync.WaitGroup
	for _, restaurantGraphApiRequestUrl := range restaurantGraphApiRequestUrls {
		restaurantGraphApiRequestUrl := restaurantGraphApiRequestUrl
		fmt.Println(restaurantGraphApiRequestUrl)
		wg.Add(1)
		go func() {
			defer wg.Done()
			graphApiResponseFromRequestUrl, err := graphApi.GetGraphApiResponseFromTimeSlot(restaurantGraphApiRequestUrl)
			if err != nil {
				restaurant.GraphApiResults.Err <- err
				return
			}

			// getAvailableTableTimesFromRestaurantRequestUrlsIntoRestaurantsChannel TODO: remember to take into consideration the kitchens closing time (can't reserve 1h before kitchen closes.)
			// getAvailableTableTimesFromRestaurantRequestUrlsIntoRestaurantsChannel TODO: capture restaurants time till kitchen and restaurant closes.
			graphApiReservationTimes := raflaamoTime.GetGraphApiReservationTimes(graphApiResponseFromRequestUrl)

			graphApiReservationTimes.GetTimeSlotsInBetweenIntervals(restaurant.GraphApiResults.AvailableTimeSlotsBuffer, raflaamoRelatedTimes.AllRaflaamoReservationTimeIntervals)
		}()
	}
	wg.Wait()
	close(restaurant.GraphApiResults.Err)
	close(restaurant.GraphApiResults.AvailableTimeSlotsBuffer)
}
