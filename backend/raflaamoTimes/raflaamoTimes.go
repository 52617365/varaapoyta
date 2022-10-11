/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package raflaamoTimes

import (
	"backend/helpers"
	"backend/raflaamoGraphApiTimes"
	"fmt"
	"regexp"
	"time"
)

type TimeAndDate struct {
	CurrentTime int64
	CurrentDate string
}

type RaflaamoTimes struct {
	TimeAndDate                               *TimeAndDate
	AllFutureRaflaamoReservationTimeIntervals []int64
}

// Returns all possible raflaamoTimes intervals that can be reserved in the raflaamo reservation page.
// 11:00, 11:15, 11:30 and so on.
func (times *RaflaamoTimes) getAllRaflaamoReservingIntervalsThatAreNotInThePast() {
	allTimes := make([]int64, 0, 96)
	currentTime := times.TimeAndDate.CurrentTime
	for hour := 0; hour < 24; hour++ {
		if hour < 10 {
			formattedTimeSlotOne := helpers.ConvertStringTimeToDesiredUnixFormat(fmt.Sprintf("0%d00", hour))
			formattedTimeSlotOne += 10800 // To match finnish timezone.
			if formattedTimeSlotOne > currentTime {
				allTimes = append(allTimes, formattedTimeSlotOne)
			}

			formattedTimeSlotTwo := helpers.ConvertStringTimeToDesiredUnixFormat(fmt.Sprintf("0%d15", hour))
			formattedTimeSlotTwo += 10800 // To match finnish timezone.
			if formattedTimeSlotTwo > currentTime {
				allTimes = append(allTimes, formattedTimeSlotTwo)
			}

			formattedTimeSlotThree := helpers.ConvertStringTimeToDesiredUnixFormat(fmt.Sprintf("0%d30", hour))
			formattedTimeSlotThree += 10800 // To match finnish timezone.
			if formattedTimeSlotThree > currentTime {
				allTimes = append(allTimes, formattedTimeSlotThree)
			}
			formattedTimeSlotFour := helpers.ConvertStringTimeToDesiredUnixFormat(fmt.Sprintf("0%d45", hour))
			formattedTimeSlotFour += 10800 // To match finnish timezone.
			if formattedTimeSlotFour > currentTime {
				allTimes = append(allTimes, formattedTimeSlotFour)
			}
		}
		if hour >= 10 {
			formattedTimeSlotOne := helpers.ConvertStringTimeToDesiredUnixFormat(fmt.Sprintf("%d00", hour))
			formattedTimeSlotOne += 10800 // To match finnish timezone.
			if formattedTimeSlotOne > currentTime {
				allTimes = append(allTimes, formattedTimeSlotOne)
			}
			formattedTimeSlotTwo := helpers.ConvertStringTimeToDesiredUnixFormat(fmt.Sprintf("%d15", hour))
			formattedTimeSlotTwo += 10800 // To match finnish timezone.
			if formattedTimeSlotTwo > currentTime {
				allTimes = append(allTimes, formattedTimeSlotTwo)
			}
			formattedTimeSlotThree := helpers.ConvertStringTimeToDesiredUnixFormat(fmt.Sprintf("%d30", hour))
			formattedTimeSlotThree += 10800 // To match finnish timezone.
			if formattedTimeSlotThree > currentTime {
				allTimes = append(allTimes, formattedTimeSlotThree)
			}
			formattedTimeSlotFour := helpers.ConvertStringTimeToDesiredUnixFormat(fmt.Sprintf("%d45", hour))
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
		CurrentTime: helpers.ConvertStringTimeToDesiredUnixFormat(matchedMinutesAndSeconds),
	}

	times.TimeAndDate = currentTimeAndDate
}

func (times *RaflaamoTimes) graphApiUnixTimeSlotIsValid(unixTimeSlot *raflaamoGraphApiTimes.CoveredTimes, restaurantClosingTimeUnix int64) bool {
	currentTimeUnix := times.TimeAndDate.CurrentTime
	if currentTimeUnix < unixTimeSlot.TimeWindowsEnd && restaurantClosingTimeUnix > unixTimeSlot.TimeWindowStart /* I'm not 100% on this logic. */ {
		return true
	}
	return false
}

// GetAllNeededRaflaamoTimes this should be called only once somewhere in the code because it's pretty expensive to construct.
func GetAllNeededRaflaamoTimes() *RaflaamoTimes {
	raflaamoTimes := RaflaamoTimes{}
	raflaamoTimes.getCurrentTimeAndDate(helpers.TimeRegex, helpers.RegexToMatchDate)
	raflaamoTimes.getAllRaflaamoReservingIntervalsThatAreNotInThePast()

	return &raflaamoTimes
}
