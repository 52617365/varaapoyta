package raflaamoTime

import (
	"backend/graphApiResponseStructure"
	"strings"
	"time"
)

type GraphApiReservationTimes struct {
	graphApiIntervalStart       int64
	graphApiIntervalStartString string
	graphApiIntervalEnd         int64
	graphApiIntervalEndString   string
}

// GetGraphApiReservationTimes TODO: we have to validate that the intervals are valid somewhere before calling the functions in this file.
func GetGraphApiReservationTimes(graphApiResponse *graphApiResponseStructure.ParsedGraphData) *GraphApiReservationTimes {
	graphApiResponseTimeIntervals := *graphApiResponse.Intervals

	graphApiTimeIntervalStart := graphApiResponseTimeIntervals[0].From
	graphApiTimeIntervalEnd := graphApiResponseTimeIntervals[0].To

	graphApiReservationTimes := GraphApiReservationTimes{graphApiIntervalStart: graphApiTimeIntervalStart, graphApiIntervalEnd: graphApiTimeIntervalEnd}

	graphApiReservationTimes.convertStartIntervalIntoString()
	graphApiReservationTimes.convertEndIntervalIntoString()

	graphApiTimeIntervalStart = graphApiReservationTimes.convertStartIntervalBackIntoDesiredUnixFormat()
	graphApiTimeIntervalEnd = graphApiReservationTimes.convertEndIntervalBackIntoDesiredUnixFormat()

	return &graphApiReservationTimes
}

// GetTimeSlotsInBetweenIntervals TODO: debug this and find what's wrong with the function.
func (graphApiReservationTimes *GraphApiReservationTimes) GetTimeSlotsInBetweenIntervals(AllRaflaamoReservationUnixTimeIntervals []int64) []string {
	timeSlotsInBetween := make([]string, 0, 50) // TODO: reserve space in advance.
	for _, raflaamoReservationUnixTimeInterval := range AllRaflaamoReservationUnixTimeIntervals {
		const oneHour = 3600 // Restaurants don't take reservations one hour before closing.
		if raflaamoReservationUnixTimeInterval > graphApiReservationTimes.graphApiIntervalStart && raflaamoReservationUnixTimeInterval <= graphApiReservationTimes.graphApiIntervalEnd {
			raflaamoReservationUnixTimeInterval += 7200 // To match timezone
			raflaamoReservationTime := ConvertUnixMillisecondsToString(raflaamoReservationUnixTimeInterval)
			timeSlotsInBetween = append(timeSlotsInBetween, raflaamoReservationTime)
		}
	}
	return timeSlotsInBetween
}

func (graphApiReservationTimes *GraphApiReservationTimes) convertStartIntervalIntoString() {
	graphApiReservationTimes.graphApiIntervalStart += 3600000 * 3 // Adding three hours into the time to match finnish timezone.
	startIntervalString := ConvertUnixMillisecondsToString(graphApiReservationTimes.graphApiIntervalStart)

	graphApiReservationTimes.graphApiIntervalStartString = startIntervalString
}

func (graphApiReservationTimes *GraphApiReservationTimes) convertEndIntervalIntoString() {
	endInterval := graphApiReservationTimes.graphApiIntervalEnd
	endInterval += 3600000 * 3 // Adding three hours into the time to match finnish timezone.
	endIntervalString := ConvertUnixMillisecondsToString(endInterval)

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

func ConvertUnixMillisecondsToString(unixTimeToConvert int64) string {
	timeInString := time.UnixMilli(unixTimeToConvert).UTC().String()

	stringTimeFromUnix := timeRegex.FindString(timeInString)

	stringTimeFromUnix = strings.Replace(stringTimeFromUnix, ":", "", -1)

	return stringTimeFromUnix
}
