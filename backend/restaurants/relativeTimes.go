/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package restaurants

import (
	"backend/helpers"
	"backend/raflaamoRestaurantsApi"
	"strconv"
)

type CalculateClosingTime struct {
	CurrentTime int64
	ClosingTime int64
}

type CalculateClosingTimeResult struct {
	RelativeHours   int
	RelativeMinutes int
}

func (restaurants *InitializeProgram) getRelativeClosingTimes(restaurant *raflaamoRestaurantsApi.ResponseFields) (*CalculateClosingTime, *CalculateClosingTime) {
	restaurantsKitchenClosingTime := restaurant.Openingtime.Kitchentime.Ranges[0].End
	currentTime := restaurants.AllNeededRaflaamoTimes.TimeAndDate.CurrentTime

	calculateTimeTillKitchenCloses := getCalculateClosingTime(currentTime, restaurantsKitchenClosingTime)
	calculateTimeTillRestaurantCloses := getCalculateClosingTime(currentTime, restaurant.Openingtime.Restauranttime.Ranges[0].End)

	return calculateTimeTillRestaurantCloses, calculateTimeTillKitchenCloses
}

func getCalculateClosingTime(currentTime int64, closingTime string) *CalculateClosingTime {
	closingTimeConvertedToUnix := helpers.ConvertStringTimeToDesiredUnixFormat(closingTime)
	return &CalculateClosingTime{CurrentTime: currentTime, ClosingTime: closingTimeConvertedToUnix}
}

func (calculation *CalculateClosingTime) CalculateRelativeTime() *CalculateClosingTimeResult {
	relativeCalculation := calculation.ClosingTime - calculation.CurrentTime

	humanReadableRelativeCalculation := helpers.ConvertUnixSecondsToString(relativeCalculation, false)

	humanReadableRelativeMinutes, _ := strconv.Atoi(humanReadableRelativeCalculation[len(humanReadableRelativeCalculation)-2:])
	humanReadableRelativeHours, _ := strconv.Atoi(humanReadableRelativeCalculation[:len(humanReadableRelativeCalculation)-2])
	return &CalculateClosingTimeResult{RelativeMinutes: humanReadableRelativeMinutes, RelativeHours: humanReadableRelativeHours}
}

func (calculation *CalculateClosingTime) relativeCalculationIsNegative(relativeCalculation int64) bool {
	if relativeCalculation < 0 {
		return true
	}
	return false
}
