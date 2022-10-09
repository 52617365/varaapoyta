/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package raflaamoTime

import (
	"backend/unixHelpers"
	"fmt"
	"regexp"
	"time"
)

// Returns all possible raflaamoTime intervals that can be reserved in the raflaamo reservation page.
// 11:00, 11:15, 11:30 and so on.
func (times *RaflaamoTimes) getAllRaflaamoReservingIntervalsThatAreNotInThePast() {
	allTimes := make([]int64, 0, 96)
	currentTime := times.TimeAndDate.CurrentTime
	for hour := 0; hour < 24; hour++ {
		if hour < 10 {
			formattedTimeSlotOne := unixHelpers.ConvertStringTimeToDesiredUnixFormat(fmt.Sprintf("0%d00", hour))
			formattedTimeSlotOne += 10800 // To match finnish timezone.
			if formattedTimeSlotOne > currentTime {
				allTimes = append(allTimes, formattedTimeSlotOne)
			}

			formattedTimeSlotTwo := unixHelpers.ConvertStringTimeToDesiredUnixFormat(fmt.Sprintf("0%d15", hour))
			formattedTimeSlotTwo += 10800 // To match finnish timezone.
			if formattedTimeSlotTwo > currentTime {
				allTimes = append(allTimes, formattedTimeSlotTwo)
			}

			formattedTimeSlotThree := unixHelpers.ConvertStringTimeToDesiredUnixFormat(fmt.Sprintf("0%d30", hour))
			formattedTimeSlotThree += 10800 // To match finnish timezone.
			if formattedTimeSlotThree > currentTime {
				allTimes = append(allTimes, formattedTimeSlotThree)
			}
			formattedTimeSlotFour := unixHelpers.ConvertStringTimeToDesiredUnixFormat(fmt.Sprintf("0%d45", hour))
			formattedTimeSlotFour += 10800 // To match finnish timezone.
			if formattedTimeSlotFour > currentTime {
				allTimes = append(allTimes, formattedTimeSlotFour)
			}
		}
		if hour >= 10 {
			formattedTimeSlotOne := unixHelpers.ConvertStringTimeToDesiredUnixFormat(fmt.Sprintf("%d00", hour))
			formattedTimeSlotOne += 10800 // To match finnish timezone.
			if formattedTimeSlotOne > currentTime {
				allTimes = append(allTimes, formattedTimeSlotOne)
			}
			formattedTimeSlotTwo := unixHelpers.ConvertStringTimeToDesiredUnixFormat(fmt.Sprintf("%d15", hour))
			formattedTimeSlotTwo += 10800 // To match finnish timezone.
			if formattedTimeSlotTwo > currentTime {
				allTimes = append(allTimes, formattedTimeSlotTwo)
			}
			formattedTimeSlotThree := unixHelpers.ConvertStringTimeToDesiredUnixFormat(fmt.Sprintf("%d30", hour))
			formattedTimeSlotThree += 10800 // To match finnish timezone.
			if formattedTimeSlotThree > currentTime {
				allTimes = append(allTimes, formattedTimeSlotThree)
			}
			formattedTimeSlotFour := unixHelpers.ConvertStringTimeToDesiredUnixFormat(fmt.Sprintf("%d45", hour))
			formattedTimeSlotFour += 10800 // To match finnish timezone.
			if formattedTimeSlotFour > currentTime {
				allTimes = append(allTimes, formattedTimeSlotFour)
			}
		}
	}
	times.AllFutureRaflaamoReservationTimeIntervals = allTimes
}

// getCurrentTimeAndDate this should be called only once.
func (times *RaflaamoTimes) getCurrentTimeAndDate(regexToMatchTime *regexp.Regexp, regexToMatchDate *regexp.Regexp) {
	dt := time.Now().String()

	matchedMinutesAndSeconds := regexToMatchTime.FindString(dt)
	currentDate := regexToMatchDate.FindString(dt)

	currentTimeAndDate := &TimeAndDate{
		CurrentDate: currentDate,
		CurrentTime: unixHelpers.ConvertStringTimeToDesiredUnixFormat(matchedMinutesAndSeconds),
	}

	times.TimeAndDate = currentTimeAndDate
}

// GetAllNeededRaflaamoTimes this should be called only once somewhere in the code because it's pretty expensive to construct.
func GetAllNeededRaflaamoTimes(regexToMatchTime *regexp.Regexp, regexToMatchDate *regexp.Regexp) *RaflaamoTimes {
	raflaamoTimes := RaflaamoTimes{}
	raflaamoTimes.getCurrentTimeAndDate(regexToMatchTime, regexToMatchDate)
	raflaamoTimes.getAllRaflaamoReservingIntervalsThatAreNotInThePast()

	return &raflaamoTimes
}
