package timeUtils

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

//type RestaurantTime struct {
//	opening int64
//	closing int64
//}

// CoveredTimes This struct contains the timeUtils you check the graph api with, and the corresponding start and end timeUtils window that the response covers.
//type CoveredTimes struct {
//	time            int64
//	timeWindowStart int64
//	timeWindowsEnd  int64
//}

//func (kitchenTime *KitchenTime) getRestaurantTimeFromKitchenTime(restaurant *responseFields) RestaurantTime {
//	// Converting restaurant_kitchen_start_time to unix, so we can compare it easily.
//	restaurantKitchenStartTime := ConvertStringTimeToUnix(restaurant.Openingtime.Kitchentime.Ranges[0].Start)
//	// We minus 1 hour from the end timeUtils because restaurants don't take reservations before that timeUtils slot.
//	// IMPORTANT: E.g. if restaurant closes at 22:00, the last possible reservation timeUtils is 21:00.
//	const oneHourUnix int64 = 3600
//	restaurantKitchenEndingTime := ConvertStringTimeToUnix(restaurant.Openingtime.Kitchentime.Ranges[0].End) - oneHourUnix
//
//	return RestaurantTime{
//		opening: restaurantKitchenStartTime,
//		closing: restaurantKitchenEndingTime,
//	}
//}

func (timeUtils *TimeUtils) getStringTimeFromCurrentTime() string {
	timeRegex, _ := regexp.Compile(`\d{2}:\d{2}`)

	timeInString := time.Unix(timeUtils.CurrentTime.CurrentTime, 0).UTC().String()

	stringTimeFromUnix := timeRegex.FindString(timeInString)

	stringTimeFromUnix = strings.Replace(stringTimeFromUnix, ":", "", -1)
	return stringTimeFromUnix
}

func (timeUtils *TimeUtils) getStringTimeFromTimeSlot() string {
	timeRegex, _ := regexp.Compile(`\d{2}:\d{2}`)

	timeInString := time.Unix(timeUtils.CurrentTime.CurrentTime, 0).UTC().String()

	stringTimeFromUnix := timeRegex.FindString(timeInString)

	stringTimeFromUnix = strings.Replace(stringTimeFromUnix, ":", "", -1)
	return stringTimeFromUnix
}

func ConvertStringTimeToUnix(timeToConvert string) int64 {
	timeToConvert = strings.Replace(timeToConvert, ":", "", -1)
	if relativeTimeFormatIsInvalid(timeToConvert) {
		return -1
	}

	minutes, _ := strconv.Atoi(timeToConvert[len(timeToConvert)-2:])
	hour, _ := strconv.Atoi(timeToConvert[:len(timeToConvert)-2])

	// if hour is 0-5 it sets day to 2 (unix)
	if hour < 5 {
		t := time.Date(1970, time.January, 2, hour, minutes, 00, 0, time.UTC)
		return t.Unix()
	}
	// if hour is 5-23 it sets day to 1 (unix)
	if hour >= 5 {
		t := time.Date(1970, time.January, 1, hour, minutes, 00, 0, time.UTC)
		return t.Unix()
	}
	return -1
}
