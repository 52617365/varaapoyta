/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package restaurants

import (
	"backend/graphApiResponseStructure"
	"backend/raflaamoGraphApi"
	"backend/raflaamoGraphApiTimes"
	"backend/raflaamoRestaurantsApi"
	"errors"
	"fmt"
)

type RestaurantWithAvailableTables struct {
	Restaurant      raflaamoRestaurantsApi.ResponseFields `json:"restaurant"`
	AvailableTables []string                              `json:"availableTables"`
}

func (initializedProgram *InitializeProgram) GetRestaurantsAndAvailableTables() ([]RestaurantWithAvailableTables, error) {
	currentTimeUnix := initializedProgram.AllNeededRaflaamoTimes.TimeAndDate.CurrentTime
	allRestaurantsFromSpecifiedCity, err := initializedProgram.RestaurantsApi.GetRestaurantsFromRaflaamoApi(currentTimeUnix)
	if err != nil {
		return nil, err
	}

	if len(allRestaurantsFromSpecifiedCity) == 0 {
		return nil, fmt.Errorf("[GetRestaurantsAndAvailableTables] - %w", errors.New("getting restaurant data succeeded but there was no data to show to the user, most likely a bug, contact the developer"))
	}

	restaurantsWithTables, err := initializedProgram.iterateRestaurants(allRestaurantsFromSpecifiedCity)
	if err != nil {
		return nil, err
	}
	return restaurantsWithTables, nil
}

func (initializedProgram *InitializeProgram) iterateRestaurants(restaurantsToIterate []raflaamoRestaurantsApi.ResponseFields) ([]RestaurantWithAvailableTables, error) {
	restaurantsWithOpenTables := make([]RestaurantWithAvailableTables, 0, 50)

	for _, restaurant := range restaurantsToIterate {
		availableTablesForRestaurant, err := initializedProgram.getAvailableTablesForRestaurant(&restaurant)
		if err != nil {
			return nil, err
		}
		restaurantWithTables := RestaurantWithAvailableTables{AvailableTables: availableTablesForRestaurant, Restaurant: restaurant}
		restaurantsWithOpenTables = append(restaurantsWithOpenTables, restaurantWithTables)
	}
	return restaurantsWithOpenTables, nil
}

func (initializedProgram *InitializeProgram) getAvailableTablesForRestaurant(restaurant *raflaamoRestaurantsApi.ResponseFields) ([]string, error) {
	raflaamoGraphApiRequestUrlStruct := raflaamoGraphApi.GetRequestUrl(restaurant.Links.TableReservationLocalized.FiFi, initializedProgram.AmountOfEaters, initializedProgram.AllNeededRaflaamoTimes.TimeAndDate.CurrentDate)
	initializedProgram.addRelativeTimesAndReservationIdToRestaurant(restaurant, raflaamoGraphApiRequestUrlStruct)

	restaurantGraphApiRequestUrls := initializedProgram.GraphApi.GenerateGraphApiRequestUrlsForRestaurant(restaurant, initializedProgram.AllNeededRaflaamoTimes.TimeAndDate.CurrentTime, initializedProgram.AllNeededRaflaamoTimes.TimeAndDate.CurrentDate, initializedProgram.AmountOfEaters)

	kitchenClosingTime := restaurant.Openingtime.Kitchentime.Ranges[0].End
	openTablesFromGraphApi, err := initializedProgram.getAvailableTableTimeSlotsFromRestaurantUrls(restaurantGraphApiRequestUrls, kitchenClosingTime)
	if err != nil { // TODO: use named errors for clarity
		if errors.As(err, &RaflaamoGraphApiDown{}) {
			// can't get any open tables, we will not continue this error.
			return nil, RaflaamoGraphApiDown{}
		}
	}

	return openTablesFromGraphApi, nil
}

type RaflaamoGraphApiDown struct {
}

func (RaflaamoGraphApiDown) Error() string {
	return "raflaamo open tables api down, we can not get open tables at this time"
}
func (initializedProgram *InitializeProgram) getAvailableTableTimeSlotsFromRestaurantUrls(restaurantGraphApiUrlTimeSlots []string, kitchenClosingTime string) ([]string, error) {
	allCapturedTimeSlots := make([]string, 0, 96)
	for _, timeSlotUrl := range restaurantGraphApiUrlTimeSlots {
		graphApiResponseFromRequestUrl, err := initializedProgram.GraphApi.GetGraphApiResponseFromTimeSlot(timeSlotUrl)
		if err != nil {
			if errors.As(err, &raflaamoGraphApi.NoAvailableTimeSlots{}) {
				continue
			}
			return nil, RaflaamoGraphApiDown{}
		}
		timeSlots := initializedProgram.captureTimeSlots(graphApiResponseFromRequestUrl, kitchenClosingTime)
		allCapturedTimeSlots = append(allCapturedTimeSlots, timeSlots...)
	}
	return allCapturedTimeSlots, nil
}

func (initializedProgram *InitializeProgram) captureTimeSlots(graphApiResponseFromRequestUrl *graphApiResponseStructure.ParsedGraphData, kitchenClosingTime string) []string {
	graphApiReservationTimes := raflaamoGraphApiTimes.GetGraphApiReservationTimes(graphApiResponseFromRequestUrl)

	timeSlotsForRestaurant := graphApiReservationTimes.GetTimeSlotsInBetweenUnixIntervals(kitchenClosingTime, initializedProgram.AllNeededRaflaamoTimes.AllFutureRaflaamoReservationTimeIntervals)

	timeSlotsCaptured := make([]string, 0, 50)
	// Capturing all the time slots for the specified time slot.
	for _, timeSlot := range timeSlotsForRestaurant {
		timeSlotsCaptured = append(timeSlotsCaptured, timeSlot)
	}
	return timeSlotsCaptured
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

	restaurant.Links.TableReservationLocalizedId = graphApiRequestUrl.IdFromReservationPageUrl // Storing the id for the front end, so we can in the future reserve with the id.
}
