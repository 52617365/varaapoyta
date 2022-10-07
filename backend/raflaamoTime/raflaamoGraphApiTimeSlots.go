package raflaamoTime

import (
	"fmt"
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
	restaurantClosingTimeUnix := ConvertStringTimeToDesiredUnixFormat(restaurantsKitchenClosingTime)
	allPossibleGraphApiTimeSlots := &[...]CoveredTimes{
		{time: 7200, timeWindowStart: 0, timeWindowsEnd: 21600},
		{time: 28800, timeWindowStart: 21600, timeWindowsEnd: 43200},
		{time: 50400, timeWindowStart: 43200, timeWindowsEnd: 64800},
		{time: 72000, timeWindowStart: 64800, timeWindowsEnd: 86400},
	}

	timeSlotsFromCurrentTimeForward := make([]string, 0, len(allPossibleGraphApiTimeSlots))
	for _, unixTimeSlot := range allPossibleGraphApiTimeSlots {
		if times.unixTimeSlotIsValid(&unixTimeSlot, restaurantClosingTimeUnix) {
			unixTimeSlotConvertedToString := unixTimeSlot.ConvertUnixTimeToString()
			timeSlotsFromCurrentTimeForward = append(timeSlotsFromCurrentTimeForward, unixTimeSlotConvertedToString)
		}
	}
	times.AllGraphApiTimeIntervals = timeSlotsFromCurrentTimeForward
}

func (times *RaflaamoTimes) unixTimeSlotIsValid(unixTimeSlot *CoveredTimes, restaurantClosingTimeUnix int64) bool {
	currentTimeUnix := times.TimeAndDate.CurrentTime
	if currentTimeUnix < unixTimeSlot.timeWindowsEnd && restaurantClosingTimeUnix > unixTimeSlot.timeWindowStart /* I'm not 100% on this logic. */ {
		return true
	}
	return false
}

// getCurrentTimeAndDate this should be called only once.
func (times *RaflaamoTimes) getCurrentTimeAndDate(regexToMatchTime *regexp.Regexp, regexToMatchDate *regexp.Regexp) {
	dt := time.Now().String()

	matchedMinutesAndSeconds := regexToMatchTime.FindString(dt)
	currentDate := regexToMatchDate.FindString(dt)

	currentTimeAndDate := &TimeAndDate{
		CurrentDate: currentDate,
		CurrentTime: ConvertStringTimeToDesiredUnixFormat(matchedMinutesAndSeconds),
	}

	times.TimeAndDate = currentTimeAndDate
}

// Returns all possible raflaamoTime intervals that can be reserved in the raflaamo reservation page.
// 11:00, 11:15, 11:30 and so on.
func (times *RaflaamoTimes) getAllRaflaamoReservingIntervalsThatAreNotInThePast() {
	allTimes := make([]int64, 0, 96)
	for hour := 0; hour < 24; hour++ {
		if hour < 10 {
			formattedTimeSlotOne := ConvertStringTimeToDesiredUnixFormat(fmt.Sprintf("0%d00", hour))
			formattedTimeSlotOne += 10800 // To match finnish timezone.
			allTimes = append(allTimes, formattedTimeSlotOne)
			formattedTimeSlotTwo := ConvertStringTimeToDesiredUnixFormat(fmt.Sprintf("0%d15", hour))
			formattedTimeSlotTwo += 10800 // To match finnish timezone.
			allTimes = append(allTimes, formattedTimeSlotTwo)
			formattedTimeSlotThree := ConvertStringTimeToDesiredUnixFormat(fmt.Sprintf("0%d30", hour))
			formattedTimeSlotThree += 10800 // To match finnish timezone.
			allTimes = append(allTimes, formattedTimeSlotThree)
			formattedTimeSlotFour := ConvertStringTimeToDesiredUnixFormat(fmt.Sprintf("0%d45", hour))
			formattedTimeSlotFour += 10800 // To match finnish timezone.
			allTimes = append(allTimes, formattedTimeSlotFour)
		}
		if hour >= 10 {
			formattedTimeSlotOne := ConvertStringTimeToDesiredUnixFormat(fmt.Sprintf("%d00", hour))
			formattedTimeSlotOne += 10800 // To match finnish timezone.
			allTimes = append(allTimes, formattedTimeSlotOne)
			formattedTimeSlotTwo := ConvertStringTimeToDesiredUnixFormat(fmt.Sprintf("%d15", hour))
			formattedTimeSlotTwo += 10800 // To match finnish timezone.
			allTimes = append(allTimes, formattedTimeSlotTwo)
			formattedTimeSlotThree := ConvertStringTimeToDesiredUnixFormat(fmt.Sprintf("%d30", hour))
			formattedTimeSlotThree += 10800 // To match finnish timezone.
			allTimes = append(allTimes, formattedTimeSlotThree)
			formattedTimeSlotFour := ConvertStringTimeToDesiredUnixFormat(fmt.Sprintf("%d45", hour))
			formattedTimeSlotFour += 10800 // To match finnish timezone.
			allTimes = append(allTimes, formattedTimeSlotFour)
		}
	}
	times.AllRaflaamoReservationTimeIntervals = allTimes
}

func (coveredTimes *CoveredTimes) ConvertUnixTimeToString() string {
	timeInString := time.Unix(coveredTimes.time, 0).UTC().String()

	stringTimeFromUnix := timeRegex.FindString(timeInString)

	stringTimeFromUnix = strings.Replace(stringTimeFromUnix, ":", "", -1)
	return stringTimeFromUnix
}

func (times *RaflaamoTimes) hourIsNotInThePast(timeSlotUnix int64) bool {
	currentTime := times.TimeAndDate.CurrentTime
	if timeSlotUnix > currentTime {
		return true
	}
	return false
}

// GetAllNeededRaflaamoTimes this should be called only once somewhere in the code because it's pretty expensive to construct.
func GetAllNeededRaflaamoTimes(regexToMatchTime *regexp.Regexp, regexToMatchDate *regexp.Regexp) *RaflaamoTimes {
	raflaamoTimes := RaflaamoTimes{}
	raflaamoTimes.getCurrentTimeAndDate(regexToMatchTime, regexToMatchDate)
	raflaamoTimes.getAllRaflaamoReservingIntervalsThatAreNotInThePast()

	return &raflaamoTimes
}