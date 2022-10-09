/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package raflaamoTime

import (
	"backend/unixHelpers"
	"regexp"
	"strings"
	"time"
)

var timeRegex = regexp.MustCompile(`\d{2}:\d{2}`)

/*
02:00 covers(00:00-06:00), 08:00 covers(6:00-12:00), 14:00 covers(12:00-18:00), 20:00 covers(18:00-00:00).
The function gets all the time windows we need to check to avoid checking redundant time windows from the past.
*/
func (times *RaflaamoTimes) GetAllGraphApiUnixTimeIntervalsFromCurrentPointForward(restaurantsKitchenClosingTime string) {
	restaurantsKitchenClosingTime = strings.ReplaceAll(restaurantsKitchenClosingTime, ":", "")
	restaurantClosingTimeUnix := unixHelpers.ConvertStringTimeToDesiredUnixFormat(restaurantsKitchenClosingTime)
	allPossibleGraphApiTimeSlots := &[...]CoveredTimes{
		{time: 7200, timeWindowStart: 0, timeWindowsEnd: 21600},
		{time: 28800, timeWindowStart: 21600, timeWindowsEnd: 43200},
		{time: 50400, timeWindowStart: 43200, timeWindowsEnd: 64800},
		{time: 72000, timeWindowStart: 64800, timeWindowsEnd: 86400},
	}

	timeSlotsFromCurrentTimeForward := make([]string, 0, len(allPossibleGraphApiTimeSlots))
	for _, graphApiUnixTimeSlot := range allPossibleGraphApiTimeSlots {
		if times.graphApiUnixTimeSlotIsValid(&graphApiUnixTimeSlot, restaurantClosingTimeUnix) {
			graphApiUnixTimeSlotIntoString := graphApiUnixTimeSlot.ConvertUnixTimeToString()
			timeSlotsFromCurrentTimeForward = append(timeSlotsFromCurrentTimeForward, graphApiUnixTimeSlotIntoString)
		}
	}
	times.AllFutureGraphApiTimeIntervals = timeSlotsFromCurrentTimeForward // TODO: does this stay?
}

func (coveredTimes *CoveredTimes) ConvertUnixTimeToString() string {
	timeInString := time.Unix(coveredTimes.time, 0).UTC().String()

	stringTimeFromUnix := timeRegex.FindString(timeInString)

	stringTimeFromUnix = strings.Replace(stringTimeFromUnix, ":", "", -1)
	return stringTimeFromUnix
}
func (times *RaflaamoTimes) graphApiUnixTimeSlotIsValid(unixTimeSlot *CoveredTimes, restaurantClosingTimeUnix int64) bool {
	currentTimeUnix := times.TimeAndDate.CurrentTime
	if currentTimeUnix < unixTimeSlot.timeWindowsEnd && restaurantClosingTimeUnix > unixTimeSlot.timeWindowStart /* I'm not 100% on this logic. */ {
		return true
	}
	return false
}
