/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package raflaamoTime

import (
	"errors"
	"fmt"
	"strconv"
)

type calculateClosingTime struct {
	currentTime int64
	closingTime string
}

func GetCalculateClosingTime(currentTime int64, closingTime string) *calculateClosingTime {
	return &calculateClosingTime{currentTime: currentTime, closingTime: closingTime}
}

type relativeHoursAndMinutes struct {
	RelativeHours   int
	RelativeMinutes int
}

func (calculation *calculateClosingTime) CalculateRelativeTime() (*relativeHoursAndMinutes, error) {
	closingTimeConvertedToUnix := ConvertStringTimeToDesiredUnixFormat(calculation.closingTime)
	relativeCalculation := closingTimeConvertedToUnix - calculation.currentTime

	if calculation.relativeCalculationIsNegative(relativeCalculation) {
		return nil, fmt.Errorf("[calculateRelative] - %w", errors.New("relative calculation was negative")) // TODO: handle.
	}

	humanReadableRelativeCalculation := ConvertUnixSecondsToString(relativeCalculation, false)

	humanReadableRelativeMinutes, _ := strconv.Atoi(humanReadableRelativeCalculation[len(humanReadableRelativeCalculation)-2:])
	humanReadableRelativeHours, _ := strconv.Atoi(humanReadableRelativeCalculation[:len(humanReadableRelativeCalculation)-2])
	return &relativeHoursAndMinutes{RelativeMinutes: humanReadableRelativeMinutes, RelativeHours: humanReadableRelativeHours}, nil
}

func (calculation *calculateClosingTime) relativeCalculationIsNegative(relativeCalculation int64) bool {
	if relativeCalculation < 0 {
		return true
	}
	return false
}
