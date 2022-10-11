/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package restaurants

import (
	"backend/raflaamoGraphApi"
	"backend/raflaamoGraphApiTimes"
	"backend/raflaamoRestaurantsApi"
	"fmt"
	"log"
)

// TODO: capture results into channel first then into a string slice.
func (initProgram *InitializeProgram) getAvailableTableTimeSlotsFromRestaurantUrls(restaurantGraphApiUrlTimeSlots []string, kitchenClosingTime string) []string {
	graphApiResults := make(chan string, len(restaurantGraphApiUrlTimeSlots))
	err := make(chan error, len(restaurantGraphApiUrlTimeSlots))

	for _, graphApiUrlTimeSlot := range restaurantGraphApiUrlTimeSlots {
		graphApiResponseFromRequestUrl, err := initProgram.GraphApi.GetGraphApiResponseFromTimeSlot(graphApiUrlTimeSlot)
		if err != nil {
			// TODO: if error has to do with not being able to access graphApi then terminate or something else.
			continue
		}
		if intervals := *graphApiResponseFromRequestUrl.Intervals; intervals[0].Color == "transparent" {
			// send something along the channel like an empty string or something else.
		} else {
			// send result along the channel
			graphApiReservationTimes := raflaamoGraphApiTimes.GetGraphApiReservationTimes(graphApiResponseFromRequestUrl)

			timeSlotsForRestaurant := graphApiReservationTimes.GetTimeSlotsInBetweenUnixIntervals(kitchenClosingTime, initProgram.AllNeededRaflaamoTimes.AllFutureRaflaamoReservationTimeIntervals)

			// Capturing all the time slots for the specified time slot.
			for _, timeSlot := range timeSlotsForRestaurant {
				graphApiResults <- timeSlot
			}
			fmt.Println(timeSlotsForRestaurant) // TODO: capture
		}
	}
}

func (initProgram *InitializeProgram) getAvailableTablesForRestaurant(restaurant *raflaamoRestaurantsApi.ResponseFields) ([]string, error) {
	raflaamoGraphApiRequestUrlStruct := raflaamoGraphApi.GetRequestUrl(restaurant.Links.TableReservationLocalized.FiFi, initProgram.AmountOfEaters, initProgram.AllNeededRaflaamoTimes.TimeAndDate.CurrentDate)
	initProgram.addRelativeTimesAndReservationIdToRestaurant(restaurant, raflaamoGraphApiRequestUrlStruct)

	restaurantGraphApiRequestUrls := initProgram.GraphApi.GenerateGraphApiRequestUrlsForRestaurant(restaurant, initProgram)

	kitchenClosingTime := restaurant.Openingtime.Kitchentime.Ranges[0].End
	openTablesFromGraphApi := initProgram.getAvailableTableTimeSlotsFromRestaurantUrls(restaurantGraphApiRequestUrls, kitchenClosingTime)

	return openTablesFromGraphApi, nil
}

func (initProgram *InitializeProgram) iterateRestaurants(restaurantsToIterate []raflaamoRestaurantsApi.ResponseFields) ([]raflaamoRestaurantsApi.ResponseFields, error) {
	restaurantsWithOpenTables := make([]raflaamoRestaurantsApi.ResponseFields, 0, 50)

	for _, restaurant := range restaurantsToIterate {
		resultsForRestaurant, err := initProgram.getAvailableTablesForRestaurant(&restaurant)
		if err != nil {
			continue
		}
		restaurant.AvailableTimeSlots = resultsForRestaurant // TODO: make sure this sticks.
		restaurantsWithOpenTables = append(restaurantsWithOpenTables, restaurant)
	}
	return restaurantsWithOpenTables, nil
}

func (initProgram *InitializeProgram) getRestaurantsAndAvailableTables() ([]raflaamoRestaurantsApi.ResponseFields, error) {
	currentTime := initProgram.AllNeededRaflaamoTimes.TimeAndDate.CurrentTime
	restaurantsFromApi, err := initProgram.RestaurantsApi.GetRestaurantsFromRaflaamoApi(currentTime)
	if err != nil {
		return nil, fmt.Errorf("server down or raflaamo down")
	}
	restaurantsWithTables, err := initProgram.iterateRestaurants(restaurantsFromApi)
	if err != nil {
		log.Fatalln("") // TODO: handle
	}
	return restaurantsWithTables, nil
	// TODO: iterate initProgram
	// TODO: getAvailableTablesForRestaurant
}

// downwards from here is old code.

func (initProgram *InitializeProgram) addRelativeTimesAndReservationIdToRestaurant(restaurant *raflaamoRestaurantsApi.ResponseFields, graphApiRequestUrl *raflaamoGraphApi.RequestUrl) {
	timeTillRestaurantCloses, timeTillKitchenCloses := initProgram.getRelativeClosingTimes(restaurant)

	restaurantRelativeTime := timeTillRestaurantCloses.CalculateRelativeTime()
	restaurant.Openingtime.TimeTillRestaurantClosedMinutes = restaurantRelativeTime.RelativeMinutes
	restaurant.Openingtime.TimeTillRestaurantClosedHours = restaurantRelativeTime.RelativeHours

	kitchenRelativeTime := timeTillKitchenCloses.CalculateRelativeTime()
	restaurant.Openingtime.TimeLeftToReserveHours = kitchenRelativeTime.RelativeHours
	restaurant.Openingtime.TimeLeftToReserveMinutes = kitchenRelativeTime.RelativeMinutes

	// TODO: this does not actually stick in the long term.
	restaurant.Links.TableReservationLocalizedId = graphApiRequestUrl.IdFromReservationPageUrl // Storing the id for the front end.
}
