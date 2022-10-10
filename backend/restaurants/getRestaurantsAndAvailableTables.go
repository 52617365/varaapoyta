/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package restaurants

import (
	"backend/raflaamoGraphApi"
	"backend/raflaamoRestaurantsApi"
	"backend/raflaamoTime"
	"backend/regex"
	"log"
	"sync"
)

func getRestaurants(city string, amountOfEaters string) *Restaurants {
	allNeededRaflaamoTimes := raflaamoTime.GetAllNeededRaflaamoTimes(regex.RegexToMatchTime, regex.RegexToMatchDate)
	graphApi := raflaamoGraphApi.GetRaflaamoGraphApi()
	initializedRaflaamoRestaurantsApi := raflaamoRestaurantsApi.GetRaflaamoRestaurantsApi(city)
	return &Restaurants{
		City:                   city,
		AmountOfEaters:         amountOfEaters,
		AllNeededRaflaamoTimes: allNeededRaflaamoTimes,
		GraphApi:               graphApi,
		RestaurantsApi:         initializedRaflaamoRestaurantsApi,
	}
}

// TODO: id from reservation url link does not stick, same with time slots.
// TODO: Get rid of all the magical shit happening for example, modifying parameter references, this will end up in a fucking nightmare, don't do it at all, I'm serious, DO NOT.
// GetRestaurantsAndAvailableTables This is the entry point to the functionality.
func (restaurants *Restaurants) getRestaurantsAndAvailableTablesIntoChannel() []raflaamoRestaurantsApi.ResponseFields {
	currentTime := restaurants.AllNeededRaflaamoTimes.TimeAndDate.CurrentTime

	restaurantsFromApi, err := restaurants.RestaurantsApi.GetRestaurantsFromRaflaamoApi(currentTime)
	if err != nil {
		log.Fatalln(err)
	}

	for _, restaurant := range restaurantsFromApi {
		restaurant := restaurant

		go func() {
			restaurants.getAvailableTablesForRestaurant(&restaurant)
		}()
	}
	return restaurantsFromApi
}

type ResponseFields = raflaamoRestaurantsApi.ResponseFields

func (restaurants *Restaurants) getAvailableTablesForRestaurant(restaurant *raflaamoRestaurantsApi.ResponseFields) {
	restaurantsKitchenClosingTime := restaurant.Openingtime.Kitchentime.Ranges[0].End
	graphApiTimeIntervalsFromTheFuture := restaurants.AllNeededRaflaamoTimes.GetAllGraphApiUnixTimeIntervalsFromCurrentPointForward(restaurantsKitchenClosingTime)
	restaurants.AllNeededRaflaamoTimes.AllFutureGraphApiTimeIntervals = graphApiTimeIntervalsFromTheFuture

	raflaamoGraphApiRequestUrlStruct := raflaamoGraphApi.GetRaflaamoGraphApiRequestUrl(restaurant.Links.TableReservationLocalized.FiFi, restaurants.AmountOfEaters, restaurants.AllNeededRaflaamoTimes.TimeAndDate.CurrentDate)

	restaurant.Links.TableReservationLocalizedId = raflaamoGraphApiRequestUrlStruct.IdFromReservationPageUrl // Storing the id for the front end.

	restaurantGraphApiRequestUrls := raflaamoGraphApiRequestUrlStruct.GenerateGraphApiRequestUrlsForRestaurant(restaurants.AllNeededRaflaamoTimes)

	restaurants.addRelativeTimesToRestaurant(restaurant, restaurantsKitchenClosingTime)

	restaurants.getAvailableTableTimesFromRestaurantRequestUrlsIntoRestaurantsChannel(restaurant, restaurantGraphApiRequestUrls)
}

func (restaurants *Restaurants) addRelativeTimesToRestaurant(restaurant *raflaamoRestaurantsApi.ResponseFields, restaurantsKitchenClosingTime string) {
	currentTime := restaurants.AllNeededRaflaamoTimes.TimeAndDate.CurrentTime

	calculateTimeTillKitchenCloses := raflaamoTime.GetCalculateClosingTime(currentTime, restaurantsKitchenClosingTime)
	restaurants.addRelativeKitchenTime(restaurant, calculateTimeTillKitchenCloses)

	calculateTimeTillRestaurantCloses := raflaamoTime.GetCalculateClosingTime(currentTime, restaurant.Openingtime.Restauranttime.Ranges[0].End)
	restaurants.addRelativeRestaurantTime(restaurant, calculateTimeTillRestaurantCloses)
}

/* Contract:
*  currentTime should not be bigger than restaurantClosingTime.
*  This contract is currently enforced in [restaurantsApi.go] (filterBadRestaurantsOut).
 */
func (restaurants *Restaurants) addRelativeRestaurantTime(restaurant *raflaamoRestaurantsApi.ResponseFields, calculateRestaurantsClosingTime *raflaamoTime.CalculateClosingTime) {
	restaurantRelativeTime := calculateRestaurantsClosingTime.CalculateRelativeTime()
	restaurant.Openingtime.TimeTillRestaurantClosedMinutes = restaurantRelativeTime.RelativeMinutes
	restaurant.Openingtime.TimeTillRestaurantClosedHours = restaurantRelativeTime.RelativeHours
}

/* Contract:
*  currentTime should not be bigger than restaurantsKitchenClosingTime.
*  This contract is currently enforced in [restaurantsApi.go] (filterBadRestaurantsOut).
 */
func (restaurants *Restaurants) addRelativeKitchenTime(restaurant *raflaamoRestaurantsApi.ResponseFields, calculateKitchenClosingTime *raflaamoTime.CalculateClosingTime) {
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
				restaurant.GraphApiResults.AvailableTimeSlotsBuffer <- ""
				return
			}

			timeIntervals := *graphApiResponseFromRequestUrl.Intervals
			if timeIntervals[0].Color == "transparent" {
				// TODO: why is this fucking up everything?
				restaurant.GraphApiResults.AvailableTimeSlotsBuffer <- ""
				restaurant.GraphApiResults.Err <- nil
				//return
			}

			graphApiReservationTimes := raflaamoTime.GetGraphApiReservationTimes(graphApiResponseFromRequestUrl)

			graphApiReservationTimes.GetTimeSlotsInBetweenUnixIntervals(restaurant, restaurants.AllNeededRaflaamoTimes.AllFutureRaflaamoReservationTimeIntervals)
		}()
	}
	wg.Wait()
	close(restaurant.GraphApiResults.Err)
	close(restaurant.GraphApiResults.AvailableTimeSlotsBuffer)
}
