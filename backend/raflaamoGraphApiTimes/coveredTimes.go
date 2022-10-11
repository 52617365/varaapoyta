/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package raflaamoGraphApiTimes

import (
	"backend/regexHelpers"
	"backend/unixHelpers"
	"strings"
	"time"
)

// CoveredTimes This struct contains the raflaamoTimes you check the graph api with, and the corresponding start and end raflaamoTimes window that the response covers.
type CoveredTimes struct {
	Time            int64
	TimeWindowStart int64
	TimeWindowsEnd  int64
}

func (coveredTimes *CoveredTimes) ConvertUnixTimeToString() string {
	timeInString := time.Unix(coveredTimes.Time, 0).UTC().String()

	stringTimeFromUnix := regexHelpers.TimeRegex.FindString(timeInString)

	stringTimeFromUnix = strings.Replace(stringTimeFromUnix, ":", "", -1)
	return stringTimeFromUnix
}

/*
02:00 covers(00:00-06:00), 08:00 covers(6:00-12:00), 14:00 covers(12:00-18:00), 20:00 covers(18:00-00:00).
The function gets all the time windows we need to check to avoid checking redundant time windows from the past.
*/
func GetAllFutureGraphApiTimeSlots(restaurantsKitchenClosingTime string) []string {
	restaurantsKitchenClosingTime = strings.ReplaceAll(restaurantsKitchenClosingTime, ":", "")
	restaurantClosingTimeUnix := unixHelpers.ConvertStringTimeToDesiredUnixFormat(restaurantsKitchenClosingTime)
	allPossibleGraphApiTimeSlots := &[...]CoveredTimes{
		{Time: 7200, TimeWindowStart: 0, TimeWindowsEnd: 21600},
		{Time: 28800, TimeWindowStart: 21600, TimeWindowsEnd: 43200},
		{Time: 50400, TimeWindowStart: 43200, TimeWindowsEnd: 64800},
		{Time: 72000, TimeWindowStart: 64800, TimeWindowsEnd: 86400},
	}

	timeSlotsFromCurrentTimeForward := make([]string, 0, len(allPossibleGraphApiTimeSlots))
	for _, graphApiUnixTimeSlot := range allPossibleGraphApiTimeSlots {
		if graphApiUnixTimeSlotIsValid(&graphApiUnixTimeSlot, restaurantClosingTimeUnix) {
			graphApiUnixTimeSlotIntoString := graphApiUnixTimeSlot.ConvertUnixTimeToString()
			timeSlotsFromCurrentTimeForward = append(timeSlotsFromCurrentTimeForward, graphApiUnixTimeSlotIntoString)
		}
	}
	return timeSlotsFromCurrentTimeForward
}
func graphApiUnixTimeSlotIsValid(graphApiUnixTimeSlot *CoveredTimes, restaurantClosingTimeUnix int64) bool {
	return restaurantClosingTimeUnix < graphApiUnixTimeSlot.Time
}
