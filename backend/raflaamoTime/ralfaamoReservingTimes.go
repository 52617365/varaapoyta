package raflaamoTime

import (
	"backend/graphApiResponseStructure"
	"strings"
	"time"
)

type GraphApiReservationTimes struct {
	graphApiIntervalStart            int64
	graphApiIntervalStartString      string
	graphApiIntervalEnd              int64
	graphApiIntervalEndString        string
	capturedPossibleReservationTimes []string
}

// GetGraphApiReservationTimes TODO: we have to validate that the intervals are valid somewhere before calling the functions in this file.
func GetGraphApiReservationTimes(graphApiResponse *graphApiResponseStructure.ParsedGraphData) *GraphApiReservationTimes {
	graphApiResponseTimeIntervals := *graphApiResponse.Intervals

	graphApiTimeIntervalStart := graphApiResponseTimeIntervals[0].From
	graphApiTimeIntervalEnd := graphApiResponseTimeIntervals[0].To

	graphApiReservationTimes := GraphApiReservationTimes{graphApiIntervalStart: graphApiTimeIntervalStart, graphApiIntervalEnd: graphApiTimeIntervalEnd}

	return &graphApiReservationTimes
}
func (graphApiReservationTimes *GraphApiReservationTimes) GetTimeSlotsInBetweenIntervals(AllRaflaamoReservationUnixTimeIntervals []int64) []string {
	graphApiReservationTimes.convertStartIntervalIntoString()
	graphApiReservationTimes.convertEndIntervalIntoString()

	graphApiStartInDesiredUnixFormat := graphApiReservationTimes.convertStartIntervalBackIntoDesiredUnixFormat()
	graphApiEndInDesiredUnixFormat := graphApiReservationTimes.convertEndIntervalBackIntoDesiredUnixFormat()

	timeSlotsInBetween := make([]string, 0, 50) // TODO: reserve space in advance.
	for _, raflaamoReservationUnixTimeInterval := range AllRaflaamoReservationUnixTimeIntervals {
		const oneHour = 3600 // Restaurants don't take reservations one hour before closing.
		if raflaamoReservationUnixTimeInterval > graphApiStartInDesiredUnixFormat && raflaamoReservationUnixTimeInterval <= graphApiEndInDesiredUnixFormat {
			raflaamoReservationUnixTimeInterval += 7200 // To match timezone
			raflaamoReservationTime := convertUnixToStringTime(raflaamoReservationUnixTimeInterval)
			timeSlotsInBetween = append(timeSlotsInBetween, raflaamoReservationTime)
		}
	}
	return timeSlotsInBetween
}

func (graphApiReservationTimes *GraphApiReservationTimes) convertStartIntervalIntoString() {
	startInterval := graphApiReservationTimes.graphApiIntervalStart
	startIntervalString := convertUnixToStringTime(startInterval)

	graphApiReservationTimes.graphApiIntervalStartString = startIntervalString
}

func (graphApiReservationTimes *GraphApiReservationTimes) convertEndIntervalIntoString() {
	endInterval := graphApiReservationTimes.graphApiIntervalEnd
	endIntervalString := convertUnixToStringTime(endInterval)

	graphApiReservationTimes.graphApiIntervalEndString = endIntervalString
}

func (graphApiReservationTimes *GraphApiReservationTimes) convertStartIntervalBackIntoDesiredUnixFormat() int64 {
	startIntervalString := graphApiReservationTimes.graphApiIntervalStartString
	startIntervalStringInDesiredUnixFormat := ConvertStringTimeToDesiredUnixFormat(startIntervalString)
	startIntervalStringInDesiredUnixFormat += 3600 // Adding one hour because of timezone stuff
	return startIntervalStringInDesiredUnixFormat
}

func (graphApiReservationTimes *GraphApiReservationTimes) convertEndIntervalBackIntoDesiredUnixFormat() int64 {
	endIntervalString := graphApiReservationTimes.graphApiIntervalEndString
	endIntervalStringInDesiredUnixFormat := ConvertStringTimeToDesiredUnixFormat(endIntervalString)
	endIntervalStringInDesiredUnixFormat += 3600 // Adding one hour because of timezone stuff
	return endIntervalStringInDesiredUnixFormat
}

func convertUnixToStringTime(unixTimeToConvert int64) string {
	timeInString := time.Unix(unixTimeToConvert, 0).UTC().String()

	stringTimeFromUnix := timeRegex.FindString(timeInString)

	stringTimeFromUnix = strings.Replace(stringTimeFromUnix, ":", "", -1)

	return stringTimeFromUnix
}
