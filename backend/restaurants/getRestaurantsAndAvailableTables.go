/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package restaurants

import (
	"backend/raflaamoGraphApi"
	"backend/raflaamoGraphApiTimes"
	"backend/raflaamoRestaurantsApi"
	"backend/raflaamoTimes"
	"log"
	"sync"
)

type ResponseFields = raflaamoRestaurantsApi.ResponseFields

func getInitializeProgram(city string, amountOfEaters string) *InitializeProgram {
	allNeededRaflaamoTimes := raflaamoTimes.GetAllNeededRaflaamoTimes()
	graphApi := raflaamoGraphApi.GetRaflaamoGraphApi()
	initializedRaflaamoRestaurantsApi := raflaamoRestaurantsApi.GetRaflaamoRestaurantsApi(city)
	return &InitializeProgram{
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
func (restaurants *InitializeProgram) getRestaurantsAndAvailableTablesIntoChannel() []raflaamoRestaurantsApi.ResponseFields {
	currentTime := restaurants.AllNeededRaflaamoTimes.TimeAndDate.CurrentTime

	restaurantsFromApi, err := restaurants.RestaurantsApi.GetRestaurantsFromRaflaamoApi(currentTime)
	if err != nil {
		log.Fatalln(err)
	}

	for _, restaurant := range restaurantsFromApi {
		restaurant := restaurant
		// TODO: declare channels here then assign them into the restaurant for clarity?
		go func() {
			restaurants.getAvailableTablesForRestaurant(&restaurant)
		}()
	}
	return restaurantsFromApi
}

// TODO: simplify this somehow.
func (restaurants *InitializeProgram) getAvailableTablesForRestaurant(restaurant *ResponseFields) {
	raflaamoGraphApiRequestUrlStruct := raflaamoGraphApi.GetRequestUrl(restaurant.Links.TableReservationLocalized.FiFi, restaurants.AmountOfEaters, restaurants.AllNeededRaflaamoTimes.TimeAndDate.CurrentDate)
	restaurants.addAdditionalInformationToRestaurant(restaurant, raflaamoGraphApiRequestUrlStruct)

	restaurantGraphApiRequestUrls := restaurants.GraphApi.GenerateGraphApiRequestUrlsForRestaurant(restaurant, raflaamoGraphApiRequestUrlStruct)
	restaurants.getAvailableTableTimesFromRestaurantRequestUrlsIntoRestaurantsChannel(restaurant, restaurantGraphApiRequestUrls)
}

func (restaurants *InitializeProgram) addAdditionalInformationToRestaurant(restaurant *raflaamoRestaurantsApi.ResponseFields, graphApiRequestUrl *raflaamoGraphApi.RequestUrl) {
	timeTillRestaurantCloses, timeTillKitchenCloses := restaurants.getRelativeClosingTimes(restaurant)

	restaurantRelativeTime := timeTillRestaurantCloses.CalculateRelativeTime()
	restaurant.Openingtime.TimeTillRestaurantClosedMinutes = restaurantRelativeTime.RelativeMinutes
	restaurant.Openingtime.TimeTillRestaurantClosedHours = restaurantRelativeTime.RelativeHours

	kitchenRelativeTime := timeTillKitchenCloses.CalculateRelativeTime()
	restaurant.Openingtime.TimeLeftToReserveHours = kitchenRelativeTime.RelativeHours
	restaurant.Openingtime.TimeLeftToReserveMinutes = kitchenRelativeTime.RelativeMinutes

	// TODO: this does not actually stick in the long term.
	restaurant.Links.TableReservationLocalizedId = graphApiRequestUrl.IdFromReservationPageUrl // Storing the id for the front end.
}

func (restaurants *InitializeProgram) getAvailableTableTimesFromRestaurantRequestUrlsIntoRestaurantsChannel(restaurant *raflaamoRestaurantsApi.ResponseFields, restaurantGraphApiRequestUrls []string) {
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

			graphApiReservationTimes := raflaamoGraphApiTimes.GetGraphApiReservationTimes(graphApiResponseFromRequestUrl)

			graphApiReservationTimes.GetTimeSlotsInBetweenUnixIntervals(restaurant, restaurants.AllNeededRaflaamoTimes.AllFutureRaflaamoReservationTimeIntervals)
		}()
	}
	wg.Wait()
	close(restaurant.GraphApiResults.Err)
	close(restaurant.GraphApiResults.AvailableTimeSlotsBuffer)
}
