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

// GetGraphApiReservationTimes TODO: we have to validate that the intervals are valid and exist somewhere before calling the functions in this file.
func GetGraphApiReservationTimes(graphApiResponse *graphApiResponseStructure.ParsedGraphData) *GraphApiReservationTimes {
	graphApiResponseTimeIntervals := *graphApiResponse.Intervals

	graphApiTimeIntervalStart := graphApiResponseTimeIntervals[0].From
	graphApiTimeIntervalEnd := graphApiResponseTimeIntervals[0].To

	graphApiReservationTimes := GraphApiReservationTimes{graphApiIntervalStart: graphApiTimeIntervalStart, graphApiIntervalEnd: graphApiTimeIntervalEnd}

	graphApiReservationTimes.convertStartIntervalIntoString()
	graphApiReservationTimes.convertEndIntervalIntoString()

	graphApiReservationTimes.graphApiIntervalStart = graphApiReservationTimes.convertStartIntervalBackIntoDesiredUnixFormat()
	graphApiReservationTimes.graphApiIntervalEnd = graphApiReservationTimes.convertEndIntervalBackIntoDesiredUnixFormat()

	return &graphApiReservationTimes
}

// GetTimeSlotsInBetweenIntervals TODO: debug this and find what's wrong with the function.
func (graphApiReservationTimes *GraphApiReservationTimes) GetTimeSlotsInBetweenIntervals(allRaflaamoReservationUnixTimeIntervals []int64) []string {
	timeSlotsInBetween := make([]string, 0, len(allRaflaamoReservationUnixTimeIntervals)) // TODO: reserve space in advance.
	for _, raflaamoReservationUnixTimeInterval := range allRaflaamoReservationUnixTimeIntervals {
		const oneHour = 3600 // Restaurants don't take reservations one hour before closing.
		if raflaamoReservationUnixTimeInterval > graphApiReservationTimes.graphApiIntervalStart && raflaamoReservationUnixTimeInterval <= graphApiReservationTimes.graphApiIntervalEnd {
			//raflaamoReservationUnixTimeInterval += 7200 // To match timezone
			raflaamoReservationTime := ConvertUnixSecondsToString(raflaamoReservationUnixTimeInterval)
			timeSlotsInBetween = append(timeSlotsInBetween, raflaamoReservationTime)
		}
	}
	return timeSlotsInBetween
}

func (graphApiReservationTimes *GraphApiReservationTimes) convertStartIntervalIntoString() {
	graphApiReservationTimes.graphApiIntervalStart += 3600 * 3 // Adding three hours into the time to match finnish timezone.
	startIntervalString := ConvertUnixSecondsToString(graphApiReservationTimes.graphApiIntervalStart)

	graphApiReservationTimes.graphApiIntervalStartString = startIntervalString
}

func (graphApiReservationTimes *GraphApiReservationTimes) convertEndIntervalIntoString() {
	graphApiReservationTimes.graphApiIntervalEnd += 3600 * 3 // Adding three hours into the time to match finnish timezone.
	// Adding three hours into the time to match finnish timezone.

	endIntervalString := ConvertUnixSecondsToString(graphApiReservationTimes.graphApiIntervalEnd)
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

func ConvertUnixSecondsToString(unixTimeToConvert int64) string {
	timeInString := time.Unix(unixTimeToConvert, 0).UTC().String()

	stringTimeFromUnix := timeRegex.FindString(timeInString)

	stringTimeFromUnix = strings.Replace(stringTimeFromUnix, ":", "", -1)

	return stringTimeFromUnix
}
