package timeUtils

import (
	"regexp"
	"strings"
	"time"
)

// CoveredTimes This struct contains the timeUtils you check the graph api with, and the corresponding start and end timeUtils window that the response covers.
type CoveredTimes struct {
	time            int64
	timeWindowStart int64
	timeWindowsEnd  int64
}

var timeRegex = regexp.MustCompile(`\d{2}:\d{2}`)

func (coveredTimes *CoveredTimes) ConvertUnixTimeToString() string {
	timeInString := time.Unix(coveredTimes.time, 0).UTC().String()

	stringTimeFromUnix := timeRegex.FindString(timeInString)

	stringTimeFromUnix = strings.Replace(stringTimeFromUnix, ":", "", -1)
	return stringTimeFromUnix
}

type AllRaflaamoReservationUnixTimeIntervals = []int64

type TimeAndDate struct {
	CurrentTime int64
	CurrentDate string
}

type RaflaamoTimes struct {
	TimeAndDate                                     *TimeAndDate
	AllRaflaamoReservationTimeIntervals             AllRaflaamoReservationUnixTimeIntervals
	AllGraphApiTimeIntervalsFromCurrentPointForward []string
}
