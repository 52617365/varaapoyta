/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package restaurants

import (
	"backend/raflaamoRestaurantsApi"

	"golang.org/x/exp/slices"
)

func removeIndexFromSlice[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

func GetRestaurantsAndCollectResults(city string, amountOfEaters string) []raflaamoRestaurantsApi.ResponseFields {
	restaurantsInstance := getRestaurants(city, amountOfEaters)
	raflaamoRestaurants := restaurantsInstance.getRestaurantsAndAvailableTablesIntoChannel()
	raflaamoRestaurants = iterateRestaurantsAndCaptureAvailableTimeSlotsFromChannel(raflaamoRestaurants)
	return raflaamoRestaurants
}

// TODO: This function right here is fucked. FIX IT ASAP.
func iterateRestaurantsAndCaptureAvailableTimeSlotsFromChannel(raflaamoRestaurants []raflaamoRestaurantsApi.ResponseFields) []raflaamoRestaurantsApi.ResponseFields {
	for index := range raflaamoRestaurants {
		restaurant := &raflaamoRestaurants[index] // Else it won't actually be a ptr to it.
		timeSlotsForRestaurant := iterateAndCaptureRestaurantTimeSlots(restaurant)

		// We want to remove restaurants that don't have available time slots.
		if len(timeSlotsForRestaurant) == 0 {
			removeIndexFromSlice(raflaamoRestaurants, index)
			continue
		}
		restaurant.AvailableTimeSlots = timeSlotsForRestaurant
		slices.Sort(restaurant.AvailableTimeSlots)
	}
	return raflaamoRestaurants
}

// iterateAndCaptureRestaurantTimeSlots captures the results from a channel because we send it over the network as JSON.
func iterateAndCaptureRestaurantTimeSlots(restaurant *raflaamoRestaurantsApi.ResponseFields) []string {
	availableTimeSlots := make([]string, 0, 50)
	for i := 0; i < len(restaurant.GraphApiResults.AvailableTimeSlotsBuffer); i++ {
		select {
		case result := <-restaurant.GraphApiResults.AvailableTimeSlotsBuffer:
			if graphApiResponseHadNoTimeSlots(result) {
				continue
			}
			availableTimeSlots = append(availableTimeSlots, result)
		case _ = <-restaurant.GraphApiResults.Err:
			continue
		}
	}
	return availableTimeSlots
}

// In other words, if graph API response had the "transparent" field set.
func graphApiResponseHadNoTimeSlots(timeSlotResult string) bool {
	// This exists because some time slots might have "transparent" field set aka no time slots found.
	// And we are forced to send something down the channel or else it will keep waiting forever expecting n items to iterate.
	return timeSlotResult == ""
}
