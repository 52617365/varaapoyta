/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package raflaamoTime

import (
	"strconv"
	"strings"
	"time"
)

func (timeUtils *TimeUtils) getStringTimeFromCurrentTime() string {
	timeInString := time.Unix(timeUtils.CurrentTime.CurrentTime, 0).UTC().String()

	stringTimeFromUnix := timeRegex.FindString(timeInString)

	stringTimeFromUnix = strings.Replace(stringTimeFromUnix, ":", "", -1)
	return stringTimeFromUnix
}

func (timeUtils *TimeUtils) getStringTimeFromTimeSlot() string {
	timeInString := time.Unix(timeUtils.CurrentTime.CurrentTime, 0).UTC().String()

	stringTimeFromUnix := timeRegex.FindString(timeInString)

	stringTimeFromUnix = strings.Replace(stringTimeFromUnix, ":", "", -1)
	return stringTimeFromUnix
}

func relativeTimeFormatIsInvalid(ourNumber string) bool {
	if _, err := strconv.ParseInt(ourNumber, 10, 64); err != nil {
		return true
	}
	if len(ourNumber) != 4 {
		return true
	}
	if ourNumber == "" {
		return true
	}
	return false
}

func ConvertStringTimeToDesiredUnixFormat(timeToConvert string) int64 {
	timeToConvert = strings.Replace(timeToConvert, ":", "", -1)
	if relativeTimeFormatIsInvalid(timeToConvert) {
		return -1
	}

	minutes, _ := strconv.Atoi(timeToConvert[len(timeToConvert)-2:])
	hour, _ := strconv.Atoi(timeToConvert[:len(timeToConvert)-2])

	// if hour is 0-5 it sets day to 2 (unix)
	if hour < 5 {
		t := time.Date(1970, time.January, 2, hour, minutes, 00, 0, time.UTC)
		unixTimeStamp := t.Unix()
		unixTimeStamp += 10800 // Adding three hours to match with the finnish timezone.
		return t.Unix()
	}
	// if hour is 5-23 it sets day to 1 (unix)
	if hour >= 5 {
		t := time.Date(1970, time.January, 1, hour, minutes, 00, 0, time.UTC)
		unixTimeStamp := t.Unix()
		unixTimeStamp += 10800 // Adding three hours to match with the finnish timezone.
		return t.Unix()
	}
	return -1
}
