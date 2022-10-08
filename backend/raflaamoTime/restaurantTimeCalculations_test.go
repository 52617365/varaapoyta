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

	relative := calculateClosingTime{currentTime: eightAM, closingTime: twoPM}
	relativeHoursAndMinutes, err := relative.calculateRelative()
	if err != nil {
		t.Errorf("[Test_calculateClosingTime_calculateRelative] - didn't expect error")
	}
	if relativeHoursAndMinutes.relativeMinutes != "0" {
		t.Errorf("[Test_calculateClosingTime_calculateRelative] - expected relative minutes to be %s but got %s", "00", relativeHoursAndMinutes.relativeMinutes)
	}
	if relativeHoursAndMinutes.relativeHours != "6" {
		t.Errorf("[Test_calculateClosingTime_calculateRelative] - expected relative hours to be %s but got %s", "06", relativeHoursAndMinutes.relativeHours)
	}
}
