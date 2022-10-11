/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package raflaamoGraphApiTimes

import (
	"backend/graphApiResponseStructure"
	"backend/helpers"
	"backend/raflaamoRestaurantsApi"
)

type GraphApiReservationTimes struct {
	graphApiIntervalStart       int64
	graphApiIntervalStartString string
	graphApiIntervalEnd         int64
	graphApiIntervalEndString   string
}

func GetGraphApiReservationTimes(graphApiResponse *graphApiResponseStructure.ParsedGraphData) *GraphApiReservationTimes {
	graphApiResponseTimeIntervals := *graphApiResponse.Intervals

	graphApiTimeIntervalStart := graphApiResponseTimeIntervals[0].From
	graphApiTimeIntervalEnd := graphApiResponseTimeIntervals[0].To

	graphApiReservationTimes := GraphApiReservationTimes{graphApiIntervalStart: graphApiTimeIntervalStart, graphApiIntervalEnd: graphApiTimeIntervalEnd}

	convertReservationTimesIntoDesiredFormat(&graphApiReservationTimes)

	return &graphApiReservationTimes
}

func convertReservationTimesIntoDesiredFormat(graphApiReservationTimes *GraphApiReservationTimes) {
	graphApiReservationTimes.convertStartUnixIntervalIntoString(true)
	graphApiReservationTimes.convertEndUnixIntervalIntoString(true)

	graphApiReservationTimes.graphApiIntervalStart = graphApiReservationTimes.convertStartUnixIntervalBackIntoDesiredUnixFormat()
	graphApiReservationTimes.graphApiIntervalEnd = graphApiReservationTimes.convertEndUnixIntervalBackIntoDesiredUnixFormat()
}

func (graphApiReservationTimes *GraphApiReservationTimes) GetTimeSlotsInBetweenUnixIntervals(restaurant *raflaamoRestaurantsApi.ResponseFields, allRaflaamoReservationUnixTimeIntervals []int64) {
	lastPossibleReservationTime := graphApiReservationTimes.getLastPossibleReservationTime(restaurant)
	for _, raflaamoReservationUnixTimeInterval := range allRaflaamoReservationUnixTimeIntervals {
		if graphApiReservationTimes.reservationUnixTimeIntervalIsValid(raflaamoReservationUnixTimeInterval, lastPossibleReservationTime) {
			raflaamoReservationTime := helpers.ConvertUnixSecondsToString(raflaamoReservationUnixTimeInterval, false)
			restaurant.GraphApiResults.AvailableTimeSlotsBuffer <- raflaamoReservationTime
		}
	}
}

func (graphApiReservationTimes *GraphApiReservationTimes) reservationUnixTimeIntervalIsValid(raflaamoReservationUnixTimeInterval int64, lastPossibleReservationTime int64) bool {
	if raflaamoReservationUnixTimeInterval > graphApiReservationTimes.graphApiIntervalStart && raflaamoReservationUnixTimeInterval <= lastPossibleReservationTime {
		return true
	}
	return false
}
func (graphApiReservationTimes *GraphApiReservationTimes) getLastPossibleReservationTime(restaurant *raflaamoRestaurantsApi.ResponseFields) int64 {
	const oneHour = 3600 // Restaurants don't take reservations one hour before closing.
	restaurantsKitchenClosingTimeUnix := helpers.ConvertStringTimeToDesiredUnixFormat(restaurant.Openingtime.Kitchentime.Ranges[0].End)
	lastPossibleReservationTime := restaurantsKitchenClosingTimeUnix - oneHour
	return lastPossibleReservationTime
}

func (graphApiReservationTimes *GraphApiReservationTimes) convertStartUnixIntervalIntoString(convertToFinnishTime bool) {
	if convertToFinnishTime {
		graphApiReservationTimes.graphApiIntervalStart += 3600000 * 3 // Adding three hours into the Time to match finnish timezone.
	}
	startIntervalString := helpers.ConvertUnixMilliSecondsToString(graphApiReservationTimes.graphApiIntervalStart)

	graphApiReservationTimes.graphApiIntervalStartString = startIntervalString
}

func (graphApiReservationTimes *GraphApiReservationTimes) convertEndUnixIntervalIntoString(convertToFinnishTime bool) {
	if convertToFinnishTime {
		graphApiReservationTimes.graphApiIntervalEnd += 3600000 * 3 // Adding three hours into the Time to match finnish timezone.
	}

	endIntervalString := helpers.ConvertUnixMilliSecondsToString(graphApiReservationTimes.graphApiIntervalEnd)
	graphApiReservationTimes.graphApiIntervalEndString = endIntervalString
}

func (graphApiReservationTimes *GraphApiReservationTimes) convertStartUnixIntervalBackIntoDesiredUnixFormat() int64 {
	startIntervalString := graphApiReservationTimes.graphApiIntervalStartString
	startIntervalStringInDesiredUnixFormat := helpers.ConvertStringTimeToDesiredUnixFormat(startIntervalString)
	return startIntervalStringInDesiredUnixFormat
}

func (graphApiReservationTimes *GraphApiReservationTimes) convertEndUnixIntervalBackIntoDesiredUnixFormat() int64 {
	endIntervalString := graphApiReservationTimes.graphApiIntervalEndString
	endIntervalStringInDesiredUnixFormat := helpers.ConvertStringTimeToDesiredUnixFormat(endIntervalString)
	return endIntervalStringInDesiredUnixFormat
}
