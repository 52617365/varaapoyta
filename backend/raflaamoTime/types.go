package raflaamoTime

// CoveredTimes This struct contains the raflaamoTime you check the graph api with, and the corresponding start and end raflaamoTime window that the response covers.
type CoveredTimes struct {
	time            int64
	timeWindowStart int64
	timeWindowsEnd  int64
}

type TimeAndDate struct {
	CurrentTime int64
	CurrentDate string
}

type RaflaamoTimes struct {
	TimeAndDate                         *TimeAndDate
	AllRaflaamoReservationTimeIntervals []int64
	AllGraphApiTimeIntervals            []string
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
