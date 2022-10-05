package timeUtils

import "strconv"

type TimeUtils struct {
	CurrentTime *TimeAndDate
	closingTime int64
	//timeLeftTillClosed int64
}

type RelativeTime struct {
	hour    int
	minutes int
}

//func (relativeTimeStruct *RelativeTime) getHourAndMinutesFromTime(relativeTime string) {
//	if relativeTimeFormatIsInvalid(relativeTime) {
//		relativeTimeStruct.hour = -1
//		relativeTimeStruct.minutes = -1
//	}
//	minutes, _ := strconv.Atoi(relativeTime[len(relativeTime)-2:])
//	hour, _ := strconv.Atoi(relativeTime[:len(relativeTime)-2])
//	relativeTimeStruct.hour = hour
//	relativeTimeStruct.minutes = minutes
//}
//
//func (timeUtils *TimeUtils) restaurantAlreadyClosed() bool {
//	return timeUtils.closingTime <= timeUtils.CurrentTime.time
//}
//
//func (timeUtils *TimeUtils) getTimeTillRestaurantClosingTime() RelativeTime {
//	// we have decremented one hour from it before because they don't take reservations in that timeUtils slot, but here we only care about if they're still open, and they are.
//	const oneHourUnix int64 = 3600
//	closingTime := timeUtils.closingTime
//	closingTime += oneHourUnix
//
//	if timeUtils.restaurantAlreadyClosed() {
//		return RelativeTime{hour: -1, minutes: -1}
//	}
//
//	relativeTimeString := timeUtils.getStringTimeFromCurrentTime()
//
//	relativeTime := RelativeTime{}
//	relativeTime.getHourAndMinutesFromTime(relativeTimeString)
//
//	return relativeTime
//}

//	func (timeUtils *TimeUtils) getTimeLeftToReserve() RelativeTime {
//		if timeUtils.restaurantAlreadyClosed() {
//			return RelativeTime{hour: -1, minutes: -1}
//		}
//
//		relativeTimeString := timeUtils.getStringTimeFromCurrentTime()
//
//		relativeTime := RelativeTime{}
//		relativeTime.getHourAndMinutesFromTime(relativeTimeString)
//
//		return relativeTime
//	}
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
