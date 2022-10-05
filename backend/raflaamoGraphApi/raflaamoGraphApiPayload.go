package raflaamoGraphApi

import (
	"backend/timeUtils"
	"fmt"
	"regexp"
)

type raflaamoGraphApiPayload struct {
	amountOfEaters           int
	timeSlot                 string
	currentDate              string
	idFromReservationPageUrl string
	regexToMatchRestaurantId *regexp.Regexp
}

func GetRaflaamoGraphApiPayload(reservationPageUrl string, amountOfEaters int, currentDate string, regexToMatchRestaurantId *regexp.Regexp) *raflaamoGraphApiPayload {
	idFromReservationPageUrl := regexToMatchRestaurantId.FindString(reservationPageUrl)
	return &raflaamoGraphApiPayload{
		amountOfEaters:           amountOfEaters,
		currentDate:              currentDate,
		idFromReservationPageUrl: idFromReservationPageUrl,
	}
}

func (graphApiPayload *raflaamoGraphApiPayload) getRequestUrl() string {
	restaurantId := graphApiPayload.idFromReservationPageUrl
	currentDate := graphApiPayload.currentDate
	timeSlot := graphApiPayload.timeSlot
	amountOfEaters := graphApiPayload.amountOfEaters

	requestUrl := fmt.Sprintf("https://s-varaukset.fi/api/recommendations/slot/%s/%s/%s/%d", restaurantId, currentDate, timeSlot, amountOfEaters)
	return requestUrl
}

type RaflaamoTimes = timeUtils.RaflaamoTimes

func (graphApiPayload *raflaamoGraphApiPayload) IterateAllPossibleTimeSlotsAndGenerateRequestUrls(raflaamoTimes *RaflaamoTimes) []string {
	requestUrls := make([]string, 0, len(raflaamoTimes.AllGraphApiTimeIntervalsFromCurrentPointForward))
	for _, graphApiTimeInterval := range raflaamoTimes.AllGraphApiTimeIntervalsFromCurrentPointForward {
		graphApiPayload.timeSlot = graphApiTimeInterval
		requestUrl := graphApiPayload.getRequestUrl()
		requestUrls = append(requestUrls, requestUrl)
	}
	return requestUrls
}
