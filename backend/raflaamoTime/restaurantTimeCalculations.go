/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package raflaamoTime

import (
	"backend/unixHelpers"
	"strconv"
)

type CalculateClosingTime struct {
	CurrentTime int64
	ClosingTime int64
}

type CalculateClosingTimeResult struct {
	RelativeHours   int
	RelativeMinutes int
}

func GetCalculateClosingTime(currentTime int64, closingTime string) *CalculateClosingTime {
	closingTimeConvertedToUnix := unixHelpers.ConvertStringTimeToDesiredUnixFormat(closingTime)
	return &CalculateClosingTime{CurrentTime: currentTime, ClosingTime: closingTimeConvertedToUnix}
}

func (calculation *CalculateClosingTime) CalculateRelativeTime() *CalculateClosingTimeResult {
	relativeCalculation := calculation.ClosingTime - calculation.CurrentTime

	humanReadableRelativeCalculation := unixHelpers.ConvertUnixSecondsToString(relativeCalculation, false)

	humanReadableRelativeMinutes, _ := strconv.Atoi(humanReadableRelativeCalculation[len(humanReadableRelativeCalculation)-2:])
	humanReadableRelativeHours, _ := strconv.Atoi(humanReadableRelativeCalculation[:len(humanReadableRelativeCalculation)-2])
	return &CalculateClosingTimeResult{RelativeMinutes: humanReadableRelativeMinutes, RelativeHours: humanReadableRelativeHours}
}

func (calculation *CalculateClosingTime) relativeCalculationIsNegative(relativeCalculation int64) bool {
	if relativeCalculation < 0 {
		return true
	}
	return false
}
