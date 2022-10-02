package raflaamoGraphApi

// import (
// 	"backend/timeUtils"
// 	"errors"
// 	"strings"
// )

// func getAvailableTables(city string, amountOfEaters int) ([]*response_fields, error) {
// 	currentTime := timeUtils.getCurrentDateAndTime()
// 	allTimeIntervals := getAllRaflaamoTimeIntervals()
// 	timeSlotsToCheck := getGraphTimeFromCurrentPointForward(currentTime.time)

// 	raflaamo, err := init_restaurants()
// 	if err != nil {
// 		return nil, errors.New("failed to connect to the raflaamo api, contact the developer")
// 	}
// 	response, err := raflaamo.get()
// 	if err != nil {
// 		return nil, errors.New("raflaamo api most likely down")
// 	}

// 	allResults := make(chan AdditionalInformationAboutRestaurant, 50)
// 	for _, restaurant := range response {
// 		if restaurantFormatIsIncorrect(city, &restaurant) {
// 			continue
// 		}
// 		restaurantAdditionalInformation := getAdditionalInformation(restaurant, len(timeSlotsToCheck))

// 		idNotFoundErr := restaurantAdditionalInformation.add()
// 		if idNotFoundErr != nil {
// 			// This branch is most likely not taken.
// 			continue
// 		}
// 		restaurantId := restaurantAdditionalInformation.restaurant.Links.TableReservationLocalizedId

// 		// Storing the result into the slice of all results and modifying the result as a reference later.
// 		allResults <- restaurantAdditionalInformation

// 		jobs := make(chan job, len(timeSlotsToCheck))

// 		// Spawning n workers, more than n is useless.
// 		for i := 0; i <= len(timeSlotsToCheck); i++ {
// 			go worker(jobs, restaurantAdditionalInformation.time_slots)
// 		}

// 		// Spawning the jobs, this is 0-4 jobs every loop iteration on the response.
// 		for i := 0; i < len(timeSlotsToCheck); i++ {
// 			ourJob := job{
// 				slot:             &timeSlotsToCheck[i],
// 				restaurantId:     restaurantId,
// 				currentTime:      currentTime,
// 				amount_of_eaters: amountOfEaters,
// 			}
// 			jobs <- ourJob
// 		}
// 		close(jobs)
// 	}

// 	restaurantsWithOpeningTimes := make([]*response_fields, 0, 50)
// 	close(allResults)
// 	for result := range allResults {
// 		restaurant := result.restaurant
// 		kitchenTimes := result.kitchen_times
// 		apiResponse := <-result.time_slots
// 		if apiResponse.err != nil {
// 			continue
// 		}
// 		for i := 0; i <= len(result.time_slots); i++ {
// 			availableTimeSlots, err := extractAvailableTimeIntervalsFromResponse(apiResponse.value, currentTime, &kitchenTimes, allTimeIntervals)
// 			if err != nil {
// 				continue
// 			}

// 			restaurant.Available_time_slots = availableTimeSlots
// 			restaurantsWithOpeningTimes = append(restaurantsWithOpeningTimes, &restaurant)
// 		}
// 	}
// 	// @Notice: If restaurant.Available_time_slots is null/nil, there are no available timeUtils slots.
// 	return restaurantsWithOpeningTimes, nil
// }

// // Function contains several conditions which we check to determine if the restaurant is ok to use.
// // We check them all inside of this one function.
// func restaurantFormatIsIncorrect(city string, restaurant *response_fields) bool {
// 	// Converting to lower so that we don't run into problems when comparing it.
// 	city = strings.ToLower(city)

// 	// The restaurant times are nil if there are no timeUtils ranges available, we want to skip those because they are useless to us without times.
// 	if restaurant.Openingtime.Restauranttime.Ranges == nil || restaurant.Openingtime.Kitchentime.Ranges == nil {
// 		return true
// 	}
// 	if restaurant.Openingtime.Kitchentime.Ranges[0].Start == "" || restaurant.Openingtime.Kitchentime.Ranges[0].End == "" {
// 		return true
// 	}
// 	if !strings.Contains(restaurant.Links.TableReservationLocalized.Fi_FI, "https://s-varaukset.fi/online/reservation/fi/") {
// 		return true
// 	}
// 	// We check that the city is the correct one since we want to filter them out depending on the passed in city.
// 	if strings.ToLower(restaurant.Address.Municipality.Fi_FI) != city {
// 		return true
// 	}
// 	restaurantOfficeHours := get_opening_and_closing_time_from_kitchen_time(restaurant)
// 	// Checking to see if the timestamps are fucked here, so we don't have to check them later.
// 	// We have already checked that the ranges exist in the previous condition (restaurant.Openingtime.Restauranttime.Ranges != nil)
// 	isInvalidDate := restaurantOfficeHours.opening >= restaurantOfficeHours.closing
// 	return isInvalidDate
// }
