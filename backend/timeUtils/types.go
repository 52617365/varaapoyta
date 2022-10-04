package timeUtils

// CoveredTimes This struct contains the timeUtils you check the graph api with, and the corresponding start and end timeUtils window that the response covers.
type CoveredTimes struct {
	time            int64
	timeWindowStart int64
	timeWindowsEnd  int64
}

type AllRaflaamoReservationUnixTimeIntervals = []int64
type RaflaamoTimes struct {
	currentTime                                     int64
	allRaflaamoReservationTimeIntervals             AllRaflaamoReservationUnixTimeIntervals
	allGraphApiTimeIntervalsFromCurrentPointForward []CoveredTimes
}
