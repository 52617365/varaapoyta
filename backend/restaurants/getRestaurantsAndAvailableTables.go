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
	"sort"
	"sync"
)

type RestaurantWithAvailableTables struct {
	Restaurant      raflaamoRestaurantsApi.ResponseFields `json:"restaurant"`
	AvailableTables []string                              `json:"availableTables"`
}

type RestaurantsChannelResults struct {
	results *RestaurantWithAvailableTables
	err     error
}

func (initializedProgram *InitializeProgram) GetRestaurantsAndAvailableTables() ([]RestaurantWithAvailableTables, error) {
	currentTimeUnix := initializedProgram.AllNeededRaflaamoTimes.TimeAndDate.CurrentTime

	restaurantsFromRaflaamoApi, err := initializedProgram.RestaurantsApi.GetRestaurantsFromRaflaamoApi(currentTimeUnix)
	if err != nil {
		return nil, err
	}

	if len(restaurantsFromRaflaamoApi) == 0 {
		return nil, fmt.Errorf("[GetRestaurantsAndAvailableTables] - %w", errors.New("getting restaurant data succeeded but there was no data to show to the user, most likely a bug, contact the developer"))
	}

	restaurantsWithTables, err := initializedProgram.iterateRestaurants(restaurantsFromRaflaamoApi)
	if err != nil {
		return nil, err
	}
	return restaurantsWithTables, nil
}

// TODO: other goroutines work now, we are syncing the goroutines for each restaurant before proceeding in getAvailableTableTimeSlotsFromRestaurantUrls, to avoid blocking in here too we should spawn the stuff in this function with goroutines.
func (initializedProgram *InitializeProgram) iterateRestaurants(restaurantsToIterate []raflaamoRestaurantsApi.ResponseFields) ([]RestaurantWithAvailableTables, error) {
	//restaurantsWithOpenTables := make([]RestaurantWithAvailableTables, 0, 30)
	restaurantsWithOpenTables := make(chan RestaurantsChannelResults, len(restaurantsToIterate))

	var wg sync.WaitGroup
	for _, restaurant := range restaurantsToIterate {
		wg.Add(1)
		restaurant := restaurant
		go func() {
			defer wg.Done()
			availableTablesForRestaurant, err := initializedProgram.getAvailableTablesForRestaurant(&restaurant)
			if err != nil {
				restaurantsWithOpenTables <- RestaurantsChannelResults{
					results: nil,
					err:     raflaamoGraphApi.RaflaamoGraphApiDown{},
				}
			}
			restaurantWithTables := RestaurantWithAvailableTables{AvailableTables: availableTablesForRestaurant, Restaurant: restaurant}
			restaurantsWithOpenTables <- RestaurantsChannelResults{
				results: &restaurantWithTables,
				err:     nil,
			}
		}()
	}
	wg.Wait()
	close(restaurantsWithOpenTables)
	return initializedProgram.syncRestaurantsWithOpenTablesChannel(restaurantsWithOpenTables)
}

func (initializedProgram *InitializeProgram) syncRestaurantsWithOpenTablesChannel(restaurantsWithOpenTables chan RestaurantsChannelResults) ([]RestaurantWithAvailableTables, error) {
	restaurantsWithOpenTablesSync := make([]RestaurantWithAvailableTables, 0, 30)
	for restaurantWithOpenTables := range restaurantsWithOpenTables {
		if restaurantWithOpenTables.err != nil {
			return nil, restaurantWithOpenTables.err
		}
		sort.Strings(restaurantWithOpenTables.results.AvailableTables)
		restaurantsWithOpenTablesSync = append(restaurantsWithOpenTablesSync, *restaurantWithOpenTables.results)
	}
	return restaurantsWithOpenTablesSync, nil
}

func (initializedProgram *InitializeProgram) getAvailableTablesForRestaurant(restaurant *raflaamoRestaurantsApi.ResponseFields) ([]string, error) {
	raflaamoGraphApiRequestUrlStruct := raflaamoGraphApi.GetRequestUrl(restaurant.Links.TableReservationLocalized.FiFi, initializedProgram.AmountOfEaters, initializedProgram.AllNeededRaflaamoTimes.TimeAndDate.CurrentDate)
	initializedProgram.addRelativeTimesAndReservationIdToRestaurant(restaurant, raflaamoGraphApiRequestUrlStruct)

	restaurantGraphApiRequestUrls := initializedProgram.GraphApi.GenerateGraphApiRequestUrlsForRestaurant(restaurant, initializedProgram.AllNeededRaflaamoTimes.TimeAndDate.CurrentTime, initializedProgram.AllNeededRaflaamoTimes.TimeAndDate.CurrentDate, initializedProgram.AmountOfEaters)

	kitchenClosingTime := restaurant.Openingtime.Kitchentime.Ranges[0].End
	openTablesFromGraphApi, err := initializedProgram.getAvailableTableTimeSlotsFromRestaurantUrls(restaurantGraphApiRequestUrls, kitchenClosingTime)
	if err != nil {
		return nil, raflaamoGraphApi.RaflaamoGraphApiDown{}
	}

	return openTablesFromGraphApi, nil
}

type GraphApiResponse struct {
	response []string
	err      error
}

func (initializedProgram *InitializeProgram) getAvailableTableTimeSlotsFromRestaurantUrls(restaurantGraphApiUrlTimeSlots []string, kitchenClosingTime string) ([]string, error) {
	channelResult := make(chan GraphApiResponse, len(restaurantGraphApiUrlTimeSlots)) // TODO: admire this solution because it's god tier
	var wg sync.WaitGroup
	for _, timeSlotUrl := range restaurantGraphApiUrlTimeSlots {
		wg.Add(1)
		timeSlotUrl := timeSlotUrl
		go func() {
			defer wg.Done()
			graphApiResponseFromRequestUrl, err := initializedProgram.GraphApi.GetGraphApiResponseFromTimeSlot(timeSlotUrl)
			if err != nil {
				if errors.As(err, &raflaamoGraphApi.NoAvailableTimeSlots{}) {
					channelResult <- GraphApiResponse{
						response: nil,
						err:      raflaamoGraphApi.NoAvailableTimeSlots{},
					}
					return
				}
				channelResult <- GraphApiResponse{
					response: nil,
					err:      raflaamoGraphApi.RaflaamoGraphApiDown{},
				}
				return
			}
			timeSlots := initializedProgram.captureTimeSlots(graphApiResponseFromRequestUrl, kitchenClosingTime)

			channelResult <- GraphApiResponse{
				response: timeSlots,
				err:      nil,
			}
		}()
	}
	wg.Wait()
	close(channelResult)

	return initializedProgram.synchronizeGraphApiChannelResults(channelResult)
}

func (initializedProgram *InitializeProgram) synchronizeGraphApiChannelResults(channelResult chan GraphApiResponse) ([]string, error) {
	allCapturedTimeSlots := make([]string, 0, 96)
	for timeSlot := range channelResult {
		if timeSlot.err != nil {
			if errors.As(timeSlot.err, &raflaamoGraphApi.NoAvailableTimeSlots{}) {
				continue
			}
			if errors.As(timeSlot.err, &raflaamoGraphApi.RaflaamoGraphApiDown{}) {
				return nil, timeSlot.err
			}
		}
		allCapturedTimeSlots = append(allCapturedTimeSlots, timeSlot.response...)
	}
	return removeDuplicate(allCapturedTimeSlots), nil
}

func (initializedProgram *InitializeProgram) captureTimeSlots(graphApiResponseFromRequestUrl *graphApiResponseStructure.ParsedGraphData, kitchenClosingTime string) []string {
	graphApiReservationTimes := raflaamoGraphApiTimes.GetGraphApiReservationTimes(graphApiResponseFromRequestUrl)

	timeSlotsForRestaurant := graphApiReservationTimes.GetTimeSlotsInBetweenUnixIntervals(kitchenClosingTime, initializedProgram.AllNeededRaflaamoTimes.AllFutureRaflaamoReservationTimeIntervals)

	return timeSlotsForRestaurant
}

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
