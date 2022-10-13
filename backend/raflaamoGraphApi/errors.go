/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package raflaamoGraphApi

type NoAvailableTimeSlots struct {
}

func (NoAvailableTimeSlots) Error() string {
	return "there were no available time slots"
}

type RaflaamoGraphApiDown struct {
}

func (RaflaamoGraphApiDown) Error() string {
	return "raflaamo open tables api down, we can not get open tables at this time"
}
