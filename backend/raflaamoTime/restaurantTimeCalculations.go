/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package raflaamoTime

import (
	"errors"
	"fmt"
)

type calculateClosingTime struct {
	currentTime int64
	closingTime string // @Notice: should mby be a string and later calculated to unix.
}

type relativeHoursAndMinutes struct {
	relativeHours   string
	relativeMinutes string
}

func (calculation *calculateClosingTime) calculateRelative() (*relativeHoursAndMinutes, error) {
	closingTimeConvertedToUnix := ConvertStringTimeToDesiredUnixFormat(calculation.closingTime)
	relativeCalculation := closingTimeConvertedToUnix - calculation.currentTime

	if calculation.relativeCalculationIsNegative(relativeCalculation) {
		return nil, fmt.Errorf("[calculateRelative] - %w", errors.New("relative calculation was negative"))
	}

	humanReadableRelativeCalculation := ConvertUnixSecondsToString(relativeCalculation, false)
	humanReadableRelativeHours := humanReadableRelativeCalculation[:2]
	humanReadableRelativeMinutes := humanReadableRelativeCalculation[2:]
	return &relativeHoursAndMinutes{relativeMinutes: humanReadableRelativeMinutes, relativeHours: humanReadableRelativeHours}, nil
}

func (calculation *calculateClosingTime) relativeCalculationIsNegative(relativeCalculation int64) bool {
	if relativeCalculation < 0 {
		return true
	}
	return false
}
