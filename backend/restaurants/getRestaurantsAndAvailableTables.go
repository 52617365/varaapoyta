/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package restaurants

import (
	"backend/raflaamoGraphApi"
	"backend/raflaamoRestaurantsApi"
	"backend/raflaamoTime"
	"errors"
	"fmt"
	"sync"
)

func GetRestaurants(city string, amountOfEaters int) (*Restaurants, error) {
	allNeededRaflaamoTimes := raflaamoTime.GetAllNeededRaflaamoTimes(RegexToMatchTime, RegexToMatchDate)
	graphApi := raflaamoGraphApi.GetRaflaamoGraphApi()
	initializedRaflaamoRestaurantsApi, err := raflaamoRestaurantsApi.GetRaflaamoRestaurantsApi(city)
	if err != nil {
		return nil, fmt.Errorf("[GetRestaurants] - %w", errors.New("error making restaurants api"))
	}
	return &Restaurants{
		City:                   city,
		AmountOfEaters:         amountOfEaters,
		AllNeededRaflaamoTimes: allNeededRaflaamoTimes,
		GraphApi:               graphApi,
		RestaurantsApi:         initializedRaflaamoRestaurantsApi,
	}, nil
}

// GetRestaurantsAndAvailableTables This is the entry point to the functionality.
func (restaurants *Restaurants) GetRestaurantsAndAvailableTables() ([]raflaamoRestaurantsApi.ResponseFields, error) {
	currentTime := restaurants.AllNeededRaflaamoTimes.TimeAndDate.CurrentTime
	allRestaurantsFromRaflaamoRestaurantsApi, err := restaurants.RestaurantsApi.GetAllRestaurantsFromRaflaamoRestaurantsApi(currentTime)
	if err != nil {
		return nil, err
	}

	for index := range allRestaurantsFromRaflaamoRestaurantsApi {
		restaurant := &allRestaurantsFromRaflaamoRestaurantsApi[index]
		go restaurants.getAvailableTablesForRestaurant(restaurant)
	}
	return allRestaurantsFromRaflaamoRestaurantsApi, nil
}

func (restaurants *Restaurants) getAvailableTablesForRestaurant(restaurant *raflaamoRestaurantsApi.ResponseFields) {
	restaurantsKitchenClosingTime := restaurant.Openingtime.Kitchentime.Ranges[0].End
	restaurants.AllNeededRaflaamoTimes.GetAllGraphApiUnixTimeIntervalsFromCurrentPointForward(restaurantsKitchenClosingTime)

	raflaamoGraphApiRequestUrlStruct := raflaamoGraphApi.GetRaflaamoGraphApiRequestUrl(restaurant.Links.TableReservationLocalized.FiFi, restaurants.AmountOfEaters, restaurants.AllNeededRaflaamoTimes.TimeAndDate.CurrentDate, regexToMatchRestaurantId)
	restaurant.Links.TableReservationLocalizedId = raflaamoGraphApiRequestUrlStruct.IdFromReservationPageUrl // Storing the id for the front end.

	restaurantGraphApiRequestUrls := raflaamoGraphApiRequestUrlStruct.GenerateGraphApiRequestUrlsForRestaurant(restaurants.AllNeededRaflaamoTimes)

	restaurants.addRelativeTimesToRestaurant(restaurant, restaurantsKitchenClosingTime)

	restaurants.getAvailableTableTimesFromRestaurantRequestUrlsIntoRestaurantsChannel(restaurant, restaurantGraphApiRequestUrls)
}

func (restaurants *Restaurants) addRelativeTimesToRestaurant(restaurant *raflaamoRestaurantsApi.ResponseFields, restaurantsKitchenClosingTime string) {
	currentTime := restaurants.AllNeededRaflaamoTimes.TimeAndDate.CurrentTime

	calculateTimeTillKitchenCloses := raflaamoTime.GetCalculateClosingTime(currentTime, restaurantsKitchenClosingTime)
	restaurants.addRelativeKitchenTime(restaurant, calculateTimeTillKitchenCloses, currentTime)

	calculateTimeTillRestaurantCloses := raflaamoTime.GetCalculateClosingTime(currentTime, restaurant.Openingtime.Restauranttime.Ranges[0].End)
	restaurants.addRelativeRestaurantTime(restaurant, calculateTimeTillRestaurantCloses, currentTime)
}

/* Contract:
*  currentTime should not be bigger than restaurantClosingTime.
*  This contract is currently enforced in [restaurantsApi.go] (filterBadRestaurantsOut).
 */
func (restaurants *Restaurants) addRelativeRestaurantTime(restaurant *raflaamoRestaurantsApi.ResponseFields, calculateRestaurantsClosingTime *raflaamoTime.CalculateClosingTime, currentTime int64) {
	restaurantRelativeTime := calculateRestaurantsClosingTime.CalculateRelativeTime()
	restaurant.Openingtime.TimeTillRestaurantClosedMinutes = restaurantRelativeTime.RelativeMinutes
	restaurant.Openingtime.TimeTillRestaurantClosedHours = restaurantRelativeTime.RelativeHours
}

/* Contract:
*  currentTime should not be bigger than restaurantsKitchenClosingTime.
*  This contract is currently enforced in [restaurantsApi.go] (filterBadRestaurantsOut).
 */
func (restaurants *Restaurants) addRelativeKitchenTime(restaurant *raflaamoRestaurantsApi.ResponseFields, calculateKitchenClosingTime *raflaamoTime.CalculateClosingTime, currentTime int64) {
	//calculateTimeTillKitchenCloses := raflaamoTime.GetCalculateClosingTime(currentTime, restaurantsKitchenClosingTime)
	kitchenRelativeTime := calculateKitchenClosingTime.CalculateRelativeTime()
	restaurant.Openingtime.TimeLeftToReserveMinutes = kitchenRelativeTime.RelativeMinutes
	restaurant.Openingtime.TimeLeftToReserveHours = kitchenRelativeTime.RelativeHours
}

func (restaurants *Restaurants) getAvailableTableTimesFromRestaurantRequestUrlsIntoRestaurantsChannel(restaurant *raflaamoRestaurantsApi.ResponseFields, restaurantGraphApiRequestUrls []string) {
	var wg sync.WaitGroup
	for _, restaurantGraphApiRequestUrl := range restaurantGraphApiRequestUrls {
		restaurantGraphApiRequestUrl := restaurantGraphApiRequestUrl
		wg.Add(1)
		go func() {
			defer wg.Done()
			graphApiResponseFromRequestUrl, err := restaurants.GraphApi.GetGraphApiResponseFromTimeSlot(restaurantGraphApiRequestUrl)
			if err != nil {
				restaurant.GraphApiResults.Err <- err
				return
			}

			// getAvailableTableTimesFromRestaurantRequestUrlsIntoRestaurantsChannel TODO: capture restaurants time till kitchen and restaurant closes.

			graphApiReservationTimes := raflaamoTime.GetGraphApiReservationTimes(graphApiResponseFromRequestUrl)

			graphApiReservationTimes.GetTimeSlotsInBetweenUnixIntervals(restaurant, restaurants.AllNeededRaflaamoTimes.AllFutureRaflaamoReservationTimeIntervals)
		}()
	}
	wg.Wait()
	close(restaurant.GraphApiResults.Err)
	close(restaurant.GraphApiResults.AvailableTimeSlotsBuffer)
}
