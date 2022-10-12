/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package restaurants

import (
	"backend/raflaamoGraphApi"
	"backend/raflaamoGraphApiTimes"
	"backend/raflaamoRestaurantsApi"
	"errors"
	"fmt"
	"log"
)

// TODO: capture results into channel first then into a string slice.
func (initializedProgram *InitializeProgram) getAvailableTableTimeSlotsFromRestaurantUrls(restaurantGraphApiUrlTimeSlots []string, kitchenClosingTime string) []string {
	graphApiResultsChan := make(chan string, len(restaurantGraphApiUrlTimeSlots))
	errChan := make(chan error, len(restaurantGraphApiUrlTimeSlots))

	for _, graphApiUrlTimeSlot := range restaurantGraphApiUrlTimeSlots {
		graphApiResponseFromRequestUrl, err := initializedProgram.GraphApi.GetGraphApiResponseFromTimeSlot(graphApiUrlTimeSlot)
		if err != nil {
			// TODO: if error has to do with not being able to access graphApi then terminate or something else.
			graphApiResultsChan <- ""
			errChan <- errors.New("graph api most likely down")
			continue
		}
		if intervals := *graphApiResponseFromRequestUrl.Intervals; intervals[0].Color == "transparent" {
			graphApiResultsChan <- ""
			errChan <- errors.New("no available tables for time slot")
			// send something along the channel like an empty string or something else.
		} else {
			// send result along the channel
			graphApiReservationTimes := raflaamoGraphApiTimes.GetGraphApiReservationTimes(graphApiResponseFromRequestUrl)

			timeSlotsForRestaurant := graphApiReservationTimes.GetTimeSlotsInBetweenUnixIntervals(kitchenClosingTime, initializedProgram.AllNeededRaflaamoTimes.AllFutureRaflaamoReservationTimeIntervals)

			// Capturing all the time slots for the specified time slot.
			for _, timeSlot := range timeSlotsForRestaurant {
				graphApiResultsChan <- timeSlot
				errChan <- nil
			}
		}
	}
	syncedTimeSlots := initializedProgram.getSyncedResultsFromChannels(restaurantGraphApiUrlTimeSlots, errChan, graphApiResultsChan)
	return syncedTimeSlots
}

func (initializedProgram *InitializeProgram) getSyncedResultsFromChannels(restaurantGraphApiUrlTimeSlots []string, errChan chan error, graphApiResultsChan chan string) []string {
	timeSlots := make([]string, 0, 96)
	for i := 0; i < len(restaurantGraphApiUrlTimeSlots); i++ {
		syncedErr := <-errChan
		syncedResult := <-graphApiResultsChan

		if syncedErr != nil {
			continue
		}
		timeSlots = append(timeSlots, syncedResult)
	}
	return timeSlots
}

func (initializedProgram *InitializeProgram) getAvailableTablesForRestaurant(restaurant *raflaamoRestaurantsApi.ResponseFields) ([]string, error) {
	raflaamoGraphApiRequestUrlStruct := raflaamoGraphApi.GetRequestUrl(restaurant.Links.TableReservationLocalized.FiFi, initializedProgram.AmountOfEaters, initializedProgram.AllNeededRaflaamoTimes.TimeAndDate.CurrentDate)
	initializedProgram.addRelativeTimesAndReservationIdToRestaurant(restaurant, raflaamoGraphApiRequestUrlStruct)

	restaurantGraphApiRequestUrls := initializedProgram.GraphApi.GenerateGraphApiRequestUrlsForRestaurant(restaurant, initializedProgram.AllNeededRaflaamoTimes.TimeAndDate.CurrentDate, initializedProgram.AmountOfEaters)

	kitchenClosingTime := restaurant.Openingtime.Kitchentime.Ranges[0].End
	openTablesFromGraphApi := initializedProgram.getAvailableTableTimeSlotsFromRestaurantUrls(restaurantGraphApiRequestUrls, kitchenClosingTime)

	return openTablesFromGraphApi, nil
}

func (initializedProgram *InitializeProgram) iterateRestaurants(restaurantsToIterate []raflaamoRestaurantsApi.ResponseFields) ([]raflaamoRestaurantsApi.ResponseFields, error) {
	restaurantsWithOpenTables := make([]raflaamoRestaurantsApi.ResponseFields, 0, 50)

	for _, restaurant := range restaurantsToIterate {
		resultsForRestaurant, err := initializedProgram.getAvailableTablesForRestaurant(&restaurant)
		if err != nil {
			continue
		}
		restaurant.AvailableTimeSlots = resultsForRestaurant // TODO: make sure this sticks.
		restaurantsWithOpenTables = append(restaurantsWithOpenTables, restaurant)
	}
	return restaurantsWithOpenTables, nil
}

func (initializedProgram *InitializeProgram) GetRestaurantsAndAvailableTables() ([]raflaamoRestaurantsApi.ResponseFields, error) {
	currentTime := initializedProgram.AllNeededRaflaamoTimes.TimeAndDate.CurrentTime
	restaurantsFromApi, err := initializedProgram.RestaurantsApi.GetRestaurantsFromRaflaamoApi(currentTime)
	if err != nil {
		return nil, fmt.Errorf("server down or raflaamo down")
	}
	restaurantsWithTables, err := initializedProgram.iterateRestaurants(restaurantsFromApi)
	if err != nil {
		log.Fatalln("") // TODO: handle
	}
	return restaurantsWithTables, nil
	// TODO: iterate initializedProgram
	// TODO: getAvailableTablesForRestaurant
}

// downwards from here is old code.

func (initializedProgram *InitializeProgram) addRelativeTimesAndReservationIdToRestaurant(restaurant *raflaamoRestaurantsApi.ResponseFields, graphApiRequestUrl *raflaamoGraphApi.RequestUrl) {
	timeTillRestaurantCloses, timeTillKitchenCloses := initializedProgram.getRelativeClosingTimes(restaurant)

	restaurantRelativeTime := timeTillRestaurantCloses.CalculateRelativeTime()
	restaurant.Openingtime.TimeTillRestaurantClosedMinutes = restaurantRelativeTime.RelativeMinutes
	restaurant.Openingtime.TimeTillRestaurantClosedHours = restaurantRelativeTime.RelativeHours

	kitchenRelativeTime := timeTillKitchenCloses.CalculateRelativeTime()
	restaurant.Openingtime.TimeLeftToReserveHours = kitchenRelativeTime.RelativeHours
	restaurant.Openingtime.TimeLeftToReserveMinutes = kitchenRelativeTime.RelativeMinutes

	// TODO: this does not actually stick in the long term.
	restaurant.Links.TableReservationLocalizedId = graphApiRequestUrl.IdFromReservationPageUrl // Storing the id for the front end.
}
