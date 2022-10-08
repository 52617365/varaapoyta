package raflaamoTime

import (
	"backend/graphApiResponseStructure"
	"backend/raflaamoRestaurantsApi"
	"strings"
	"time"
)

type GraphApiReservationTimes struct {
	graphApiIntervalStart       int64
	graphApiIntervalStartString string
	graphApiIntervalEnd         int64
	graphApiIntervalEndString   string
}

// GetGraphApiReservationTimes TODO: we have to validate that the intervals are valid and exist somewhere before calling the functions in this file.
func GetGraphApiReservationTimes(graphApiResponse *graphApiResponseStructure.ParsedGraphData) *GraphApiReservationTimes {
	graphApiResponseTimeIntervals := *graphApiResponse.Intervals

	graphApiTimeIntervalStart := graphApiResponseTimeIntervals[0].From
	graphApiTimeIntervalEnd := graphApiResponseTimeIntervals[0].To

	graphApiReservationTimes := GraphApiReservationTimes{graphApiIntervalStart: graphApiTimeIntervalStart, graphApiIntervalEnd: graphApiTimeIntervalEnd}

	convertReservationTimesIntoDesiredFormat(&graphApiReservationTimes)

	return &graphApiReservationTimes
}

func convertReservationTimesIntoDesiredFormat(graphApiReservationTimes *GraphApiReservationTimes) {
	graphApiReservationTimes.convertStartIntervalIntoString(true)
	graphApiReservationTimes.convertEndIntervalIntoString(true)

	graphApiReservationTimes.graphApiIntervalStart = graphApiReservationTimes.convertStartIntervalBackIntoDesiredUnixFormat()
	graphApiReservationTimes.graphApiIntervalEnd = graphApiReservationTimes.convertEndIntervalBackIntoDesiredUnixFormat()
}

func (graphApiReservationTimes *GraphApiReservationTimes) GetTimeSlotsInBetweenIntervals(restaurant *raflaamoRestaurantsApi.ResponseFields, allRaflaamoReservationUnixTimeIntervals []int64) {
	lastPossibleReservationTime := graphApiReservationTimes.getLastPossibleReservationTime(restaurant)
	for _, raflaamoReservationUnixTimeInterval := range allRaflaamoReservationUnixTimeIntervals {
		if graphApiReservationTimes.reservationUnixTimeIntervalIsValid(raflaamoReservationUnixTimeInterval, lastPossibleReservationTime) {
			raflaamoReservationTime := ConvertUnixSecondsToString(raflaamoReservationUnixTimeInterval, false)
			restaurant.GraphApiResults.AvailableTimeSlotsBuffer <- raflaamoReservationTime
		}
	}
}

// reservationUnixTimeIntervalIsValid TODO: debug this and find what's wrong with the function.
func (graphApiReservationTimes *GraphApiReservationTimes) reservationUnixTimeIntervalIsValid(raflaamoReservationUnixTimeInterval int64, lastPossibleReservationTime int64) bool {
	if raflaamoReservationUnixTimeInterval > graphApiReservationTimes.graphApiIntervalStart && raflaamoReservationUnixTimeInterval <= lastPossibleReservationTime { // TODO: this logic is wrong
		return true
	}
	return false
}
func (graphApiReservationTimes *GraphApiReservationTimes) getLastPossibleReservationTime(restaurant *raflaamoRestaurantsApi.ResponseFields) int64 {
	const oneHour = 3600 // Restaurants don't take reservations one hour before closing.
	restaurantsKitchenClosingTimeUnix := ConvertStringTimeToDesiredUnixFormat(restaurant.Openingtime.Kitchentime.Ranges[0].End)
	lastPossibleReservationTime := restaurantsKitchenClosingTimeUnix - oneHour
	return lastPossibleReservationTime
}

func (graphApiReservationTimes *GraphApiReservationTimes) convertStartIntervalIntoString(convertToFinnishTime bool) {
	if convertToFinnishTime {
		graphApiReservationTimes.graphApiIntervalStart += 3600000 * 3 // Adding three hours into the time to match finnish timezone.
	}
	startIntervalString := ConvertUnixMilliSecondsToString(graphApiReservationTimes.graphApiIntervalStart)

	graphApiReservationTimes.graphApiIntervalStartString = startIntervalString
}

func (graphApiReservationTimes *GraphApiReservationTimes) convertEndIntervalIntoString(convertToFinnishTime bool) {
	if convertToFinnishTime {
		graphApiReservationTimes.graphApiIntervalEnd += 3600000 * 3 // Adding three hours into the time to match finnish timezone.
	}

	endIntervalString := ConvertUnixMilliSecondsToString(graphApiReservationTimes.graphApiIntervalEnd)
	graphApiReservationTimes.graphApiIntervalEndString = endIntervalString
}

func (graphApiReservationTimes *GraphApiReservationTimes) convertStartIntervalBackIntoDesiredUnixFormat() int64 {
	startIntervalString := graphApiReservationTimes.graphApiIntervalStartString
	startIntervalStringInDesiredUnixFormat := ConvertStringTimeToDesiredUnixFormat(startIntervalString)
	return startIntervalStringInDesiredUnixFormat
}

func (graphApiReservationTimes *GraphApiReservationTimes) convertEndIntervalBackIntoDesiredUnixFormat() int64 {
	endIntervalString := graphApiReservationTimes.graphApiIntervalEndString
	endIntervalStringInDesiredUnixFormat := ConvertStringTimeToDesiredUnixFormat(endIntervalString)
	return endIntervalStringInDesiredUnixFormat
}

func ConvertUnixSecondsToString(unixTimeToConvert int64, convertToFinnishTimezone bool) string {
	if convertToFinnishTimezone {
		unixTimeToConvert += 3 * 3600 // adding 3 hours to match finnish timezone.
	}
	timeInString := time.Unix(unixTimeToConvert, 0).UTC().String()

	stringTimeFromUnix := timeRegex.FindString(timeInString)

	stringTimeFromUnix = strings.Replace(stringTimeFromUnix, ":", "", -1)

	return stringTimeFromUnix
}
func ConvertUnixMilliSecondsToString(unixTimeToConvert int64) string {
	timeInString := time.UnixMilli(unixTimeToConvert).UTC().String()

	stringTimeFromUnix := timeRegex.FindString(timeInString)

	stringTimeFromUnix = strings.Replace(stringTimeFromUnix, ":", "", -1)

	return stringTimeFromUnix
}
