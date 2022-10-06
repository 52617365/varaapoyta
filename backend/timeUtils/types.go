package timeUtils

// CoveredTimes This struct contains the timeUtils you check the graph api with, and the corresponding start and end timeUtils window that the response covers.
type CoveredTimes struct {
	time            int64
	timeWindowStart int64
	timeWindowsEnd  int64
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
type TimeUtils struct {
	CurrentTime *TimeAndDate
	closingTime int64
	//timeLeftTillClosed int64
}

type RelativeTime struct {
	hour    int
	minutes int
}
