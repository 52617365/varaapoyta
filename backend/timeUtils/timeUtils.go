package timeUtils

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

type RestaurantTime struct {
	opening int64
	closing int64
}

type DateAndTime struct {
	date string
	time int64
}

// CoveredTimes This struct contains the timeUtils you check the graph api with, and the corresponding start and end timeUtils window that the response covers.
type CoveredTimes struct {
	time            int64
	timeWindowStart int64
	timeWindowsEnd  int64
}

func (kitchenTime *KitchenTime) getRestaurantTimeFromKitchenTime(restaurant *responseFields) RestaurantTime {
	// Converting restaurant_kitchen_start_time to unix, so we can compare it easily.
	restaurantKitchenStartTime := ConvertStringTimeToUnix(restaurant.Openingtime.Kitchentime.Ranges[0].Start)
	// We minus 1 hour from the end timeUtils because restaurants don't take reservations before that timeUtils slot.
	// IMPORTANT: E.g. if restaurant closes at 22:00, the last possible reservation timeUtils is 21:00.
	const oneHourUnix int64 = 3600
	restaurantKitchenEndingTime := ConvertStringTimeToUnix(restaurant.Openingtime.Kitchentime.Ranges[0].End) - oneHourUnix

	return RestaurantTime{
		opening: restaurantKitchenStartTime,
		closing: restaurantKitchenEndingTime,
	}
}

func (timeUtils *TimeUtils) GetStringTimeFromUnix() string {
	timeRegex, _ := regexp.Compile(`\d{2}:\d{2}`)

	rawString := time.Unix(timeUtils.timeLeftTillClosed, 0).UTC().String()

	stringTimeFromUnix := timeRegex.FindString(rawString)

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

// Gets the current timeUtils and date and initializes a struct with it.
func getDateAndTime() *DateAndTime {
	dateRegex := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	timeRegex := regexp.MustCompile(`\d{2}:\d{2}`)

	dt := time.Now().String()
	dateToString := dateRegex.FindString(dt)
	timeToString := timeRegex.FindString(dt)

	return &DateAndTime{
		date: dateToString,
		time: ConvertStringTimeToUnix(timeToString),
	}
}

/*
02:00 covers(00:00-06:00), 08:00 covers(6:00-12:00), 14:00 covers(12:00-18:00), 20:00 covers(18:00-00:00).
The function gets all the time windows we need to check to avoid checking redundant time windows from the past.
*/
func (dt *DateAndTime) getGraphTimeSlotsFromCurrentPointForward(currentTime int64) []CoveredTimes {
	// Getting current_time, so we can avoid checking times from the past.
	allPossibleUnixTimeSlots := [...]CoveredTimes{
		{time: 7200, timeWindowStart: 0, timeWindowsEnd: 21600},
		{time: 28800, timeWindowStart: 21600, timeWindowsEnd: 43200},
		{time: 50400, timeWindowStart: 43200, timeWindowsEnd: 64800},
		{time: 72000, timeWindowStart: 64800, timeWindowsEnd: 86400},
	}
	unixTimeSlotsWeWant := make([]CoveredTimes, 0, len(allPossibleUnixTimeSlots))
	for _, unixTimeSlot := range allPossibleUnixTimeSlots {
		if currentTime < unixTimeSlot.timeWindowsEnd {
			unixTimeSlotsWeWant = append(unixTimeSlotsWeWant, unixTimeSlot)
		}
	}
	return unixTimeSlotsWeWant
}
