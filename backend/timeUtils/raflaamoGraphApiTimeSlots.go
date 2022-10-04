package timeUtils

import (
	"fmt"
	"regexp"
	"time"
)

/*
02:00 covers(00:00-06:00), 08:00 covers(6:00-12:00), 14:00 covers(12:00-18:00), 20:00 covers(18:00-00:00).
The function gets all the time windows we need to check to avoid checking redundant time windows from the past.
*/
func (times *RaflaamoTimes) getAllGraphApiTimeIntervalsFromCurrentPointForward() {
	currentTimeUnix := times.timeAndDate.currentTime
	allPossibleGraphApiTimeSlots := &[...]CoveredTimes{
		{time: 7200, timeWindowStart: 0, timeWindowsEnd: 21600},
		{time: 28800, timeWindowStart: 21600, timeWindowsEnd: 43200},
		{time: 50400, timeWindowStart: 43200, timeWindowsEnd: 64800},
		{time: 72000, timeWindowStart: 64800, timeWindowsEnd: 86400},
	}

	timeSlotsFromCurrentTimeForward := make([]CoveredTimes, 0, len(allPossibleGraphApiTimeSlots))
	for _, unixTimeSlot := range allPossibleGraphApiTimeSlots {
		if currentTimeUnix < unixTimeSlot.timeWindowsEnd {
			timeSlotsFromCurrentTimeForward = append(timeSlotsFromCurrentTimeForward, unixTimeSlot)
		}
	}
	times.allGraphApiTimeIntervalsFromCurrentPointForward = timeSlotsFromCurrentTimeForward
}

// getCurrentTimeAndDate this should be called only once.
func (times *RaflaamoTimes) getCurrentTimeAndDate() {
	regexToMatchTime := regexp.MustCompile(`\d{2}:\d{2}`)
	regexToMatchDate := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	dt := time.Now().String()

	matchedMinutesAndSeconds := regexToMatchTime.FindString(dt)
	currentDate := regexToMatchDate.FindString(dt)

	currentTimeAndDate := &TimeAndDate{
		currentDate: currentDate,
		currentTime: ConvertStringTimeToUnix(matchedMinutesAndSeconds),
	}

	times.timeAndDate = currentTimeAndDate
}

// Returns all possible timeUtils intervals that can be reserved in the raflaamo reservation page.
// 11:00, 11:15, 11:30 and so on.
func (times *RaflaamoTimes) getAllRaflaamoReservingIntervalsThatAreNotInThePastOrFuture() {
	allTimes := make([]int64, 0, 96)
	for hour := 0; hour < 24; hour++ {
		if hour < 10 {
			formattedTimeSlotOne := ConvertStringTimeToUnix(fmt.Sprintf("0%d00", hour))
			if times.hourIsNotInThePast(formattedTimeSlotOne) {
				allTimes = append(allTimes, formattedTimeSlotOne)
			}
			formattedTimeSlotTwo := ConvertStringTimeToUnix(fmt.Sprintf("0%d15", hour))
			if times.hourIsNotInThePast(formattedTimeSlotTwo) {
				allTimes = append(allTimes, formattedTimeSlotTwo)
			}
			formattedTimeSlotThree := ConvertStringTimeToUnix(fmt.Sprintf("0%d30", hour))
			if times.hourIsNotInThePast(formattedTimeSlotThree) {
				allTimes = append(allTimes, formattedTimeSlotThree)
			}

			formattedTimeSlotFour := ConvertStringTimeToUnix(fmt.Sprintf("0%d45", hour))
			if times.hourIsNotInThePast(formattedTimeSlotFour) {
				allTimes = append(allTimes, formattedTimeSlotFour)
			}
		}
		if hour >= 10 {
			formattedTimeSlotOne := ConvertStringTimeToUnix(fmt.Sprintf("%d00", hour))
			if times.hourIsNotInThePast(formattedTimeSlotOne) {
				allTimes = append(allTimes, formattedTimeSlotOne)
			}
			formattedTimeSlotTwo := ConvertStringTimeToUnix(fmt.Sprintf("%d15", hour))
			if times.hourIsNotInThePast(formattedTimeSlotTwo) {
				allTimes = append(allTimes, formattedTimeSlotTwo)
			}
			formattedTimeSlotThree := ConvertStringTimeToUnix(fmt.Sprintf("%d30", hour))
			if times.hourIsNotInThePast(formattedTimeSlotThree) {
				allTimes = append(allTimes, formattedTimeSlotThree)
			}
			formattedTimeSlotFour := ConvertStringTimeToUnix(fmt.Sprintf("%d45", hour))
			if times.hourIsNotInThePast(formattedTimeSlotFour) {
				allTimes = append(allTimes, formattedTimeSlotFour)
			}
		}
		times.allRaflaamoReservationTimeIntervals = allTimes
	}
}

func (times *RaflaamoTimes) hourIsNotInThePast(timeSlotUnix int64) bool {
	currentTime := times.timeAndDate.currentTime
	if timeSlotUnix > currentTime {
		return true
	}
	return false
}

// GetRaflaamoTimes this should be called only once somewhere in the code because it's pretty expensive to construct.
func GetRaflaamoTimes() *RaflaamoTimes {
	raflaamoTimes := RaflaamoTimes{}
	raflaamoTimes.getCurrentTimeAndDate()
	raflaamoTimes.getAllRaflaamoReservingIntervalsThatAreNotInThePastOrFuture()
	raflaamoTimes.getAllGraphApiTimeIntervalsFromCurrentPointForward()

	return &raflaamoTimes
}
