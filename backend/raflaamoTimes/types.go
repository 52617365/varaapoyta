/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package raflaamoTimes

type TimeAndDate struct {
	CurrentTime int64
	CurrentDate string
}

type RaflaamoTimes struct {
	TimeAndDate                               *TimeAndDate
	AllFutureRaflaamoReservationTimeIntervals []int64
	//AllFutureGraphApiTimeIntervals            []string
}

type RelativeTime struct {
	hour    int
	minutes int
}
