/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package raflaamoTime

import (
	"strings"
	"time"
)

type TimeUtils struct {
	CurrentTime *TimeAndDate
	closingTime int64
}

func (timeUtils *TimeUtils) GetStringTimeFromCurrentTime() string {
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
