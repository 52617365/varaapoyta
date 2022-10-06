package raflaamoGraphApi

import (
	"backend/raflaamoTime"
	"fmt"
	"regexp"
)

type raflaamoGraphApiRequestUrl struct {
	amountOfEaters           int
	timeSlotToCheck          string
	currentDate              string
	idFromReservationPageUrl string
}

func GetRaflaamoGraphApiRequestUrl(reservationPageUrl string, amountOfEaters int, currentDate string, regexToMatchRestaurantId *regexp.Regexp) *raflaamoGraphApiRequestUrl {
	idFromReservationPageUrl := regexToMatchRestaurantId.FindString(reservationPageUrl)
	return &raflaamoGraphApiRequestUrl{
		amountOfEaters:           amountOfEaters,
		currentDate:              currentDate,
		idFromReservationPageUrl: idFromReservationPageUrl,
	}
}

func (graphApiPayload *raflaamoGraphApiRequestUrl) getRequestUrlForGraphApi() string {
	restaurantId := graphApiPayload.idFromReservationPageUrl
	currentDate := graphApiPayload.currentDate
	timeSlotToCheck := graphApiPayload.timeSlotToCheck
	amountOfEaters := graphApiPayload.amountOfEaters

	requestUrl := fmt.Sprintf("https://s-varaukset.fi/api/recommendations/slot/%s/%s/%s/%d", restaurantId, currentDate, timeSlotToCheck, amountOfEaters)
	return requestUrl
}

func (graphApiPayload *raflaamoGraphApiRequestUrl) GenerateGraphApiRequestUrlsForRestaurant(raflaamoTimes *raflaamoTime.RaflaamoTimes) []string {
	requestUrls := make([]string, 0, len(raflaamoTimes.AllGraphApiTimeIntervalsFromCurrentPointForward))
	for _, graphApiTimeInterval := range raflaamoTimes.AllGraphApiTimeIntervalsFromCurrentPointForward {
		graphApiPayload.timeSlotToCheck = graphApiTimeInterval
		graphApiRequestUrl := graphApiPayload.getRequestUrlForGraphApi()
		requestUrls = append(requestUrls, graphApiRequestUrl)
	}
	return requestUrls
}
