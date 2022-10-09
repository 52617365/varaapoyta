/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package raflaamoTime

import (
	"testing"
)

func Test_calculateClosingTime_calculateRelative(t *testing.T) {
	const eightAM int64 = 28800
	const twoPM string = "1400"

	relative := CalculateClosingTime{CurrentTime: eightAM, ClosingTime: twoPM}
	relativeHoursAndMinutes, err := relative.CalculateRelativeTime()
	if err != nil {
		t.Errorf("[Test_calculateClosingTime_calculateRelative] - didn't expect error")
	}
	if relativeHoursAndMinutes.RelativeMinutes != 0 {
		t.Errorf("[Test_calculateClosingTime_calculateRelative] - expected relative minutes to be %d but got %d", 0, relativeHoursAndMinutes.RelativeMinutes)
	}
	if relativeHoursAndMinutes.RelativeHours != 6 {
		t.Errorf("[Test_calculateClosingTime_calculateRelative] - expected relative hours to be %d but got %d", 6, relativeHoursAndMinutes.RelativeHours)
	}
}
