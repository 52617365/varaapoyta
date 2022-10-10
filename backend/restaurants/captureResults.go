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
	raflaamoRestaurants := restaurantsInstance.getRestaurantsAndAvailableTables()
	for index := range raflaamoRestaurants {
		restaurant := &raflaamoRestaurants[index] // or else it's not an actual ptr.

		timeSlotsForRestaurants := captureTimeSlotsForRestaurant(restaurant)
		if len(timeSlotsForRestaurants) == 0 { // We want to remove restaurants that don't have available time slots.
			removeIndexFromSlice(raflaamoRestaurants, index)
			continue
		}
		restaurant.AvailableTimeSlots = timeSlotsForRestaurants
	}
	return raflaamoRestaurants
}

func captureTimeSlotsForRestaurant(raflaamoRestaurant *raflaamoRestaurantsApi.ResponseFields) []string {
	err := <-raflaamoRestaurant.GraphApiResults.Err
	if err != nil {
		return nil
	}
	timeSlotsForRestaurant := iterateAndCaptureRestaurantTimeSlots(raflaamoRestaurant)
	slices.Sort(timeSlotsForRestaurant)
	return timeSlotsForRestaurant
}

// iterateAndCaptureRestaurantTimeSlots captures the results from a channel because we send it over the network as JSON.
func iterateAndCaptureRestaurantTimeSlots(restaurant *raflaamoRestaurantsApi.ResponseFields) []string {
	availableTimeSlots := make([]string, 0, 50)
	for availableTimeSlot := range restaurant.GraphApiResults.AvailableTimeSlotsBuffer {
		if GraphApiResponseHadNoTimeSlots(availableTimeSlot) {
			continue
		}
		availableTimeSlots = append(availableTimeSlots, availableTimeSlot)
	}
	return availableTimeSlots
}

// GraphApiResponseHadNoTimeSlots if graph API response had the "transparent" field set.
func GraphApiResponseHadNoTimeSlots(timeSlotResult string) bool {
	// This exists because some time slots might have "transparent" field set aka no time slots found.
	// And we are forced to send something down the channel or else it will keep waiting forever expecting n items to iterate.
	return timeSlotResult == ""
}
