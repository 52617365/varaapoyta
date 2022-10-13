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
	"golang.org/x/exp/slices"
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
	if err != nil {
		return nil, err
	}

	return openTablesFromGraphApi, nil
}

func (initializedProgram *InitializeProgram) getAvailableTableTimeSlotsFromRestaurantUrls(restaurantGraphApiUrlTimeSlots []string, kitchenClosingTime string) ([]string, error) {
	allCapturedTimeSlots := make([]string, 0, 96)
	for _, timeSlotUrl := range restaurantGraphApiUrlTimeSlots {
		graphApiResponseFromRequestUrl, err := initializedProgram.GraphApi.GetGraphApiResponseFromTimeSlot(timeSlotUrl)
		if err != nil {
			if errors.As(err, &raflaamoGraphApi.NoAvailableTimeSlots{}) {
				continue
			}
			return nil, fmt.Errorf("[getAvailableTableTimeSlotsFromRestaurantUrls] - %w", errors.New("raflaamo open tables api seems to be down"))
		}
		timeSlots := initializedProgram.captureTimeSlots(graphApiResponseFromRequestUrl, kitchenClosingTime)

		// avoid duplicates.
		for _, timeSlot := range timeSlots {
			if !slices.Contains(allCapturedTimeSlots, timeSlot) {
				allCapturedTimeSlots = append(allCapturedTimeSlots, timeSlot)
			}
		}
	}
	return allCapturedTimeSlots, nil
}

func (initializedProgram *InitializeProgram) captureTimeSlots(graphApiResponseFromRequestUrl *graphApiResponseStructure.ParsedGraphData, kitchenClosingTime string) []string {
	graphApiReservationTimes := raflaamoGraphApiTimes.GetGraphApiReservationTimes(graphApiResponseFromRequestUrl)

	timeSlotsForRestaurant := graphApiReservationTimes.GetTimeSlotsInBetweenUnixIntervals(kitchenClosingTime, initializedProgram.AllNeededRaflaamoTimes.AllFutureRaflaamoReservationTimeIntervals)

	return timeSlotsForRestaurant
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

func removeDuplicate[T string | int](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
