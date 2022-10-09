/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package raflaamoGraphApi

import (
	"backend/raflaamoTime"
	"fmt"
	"regexp"
)

type raflaamoGraphApiRequestUrl struct {
	amountOfEaters           string
	timeSlotToCheck          string
	currentDate              string
	IdFromReservationPageUrl string
}

func GetRaflaamoGraphApiRequestUrl(reservationPageUrl string, amountOfEaters string, currentDate string, regexToMatchRestaurantId *regexp.Regexp) *raflaamoGraphApiRequestUrl {
	idFromReservationPageUrl := regexToMatchRestaurantId.FindString(reservationPageUrl)
	return &raflaamoGraphApiRequestUrl{
		amountOfEaters:           amountOfEaters,
		currentDate:              currentDate,
		IdFromReservationPageUrl: idFromReservationPageUrl,
	}
}

func (graphApiPayload *raflaamoGraphApiRequestUrl) getRequestUrlForGraphApi() string {
	restaurantId := graphApiPayload.IdFromReservationPageUrl
	currentDate := graphApiPayload.currentDate
	timeSlotToCheck := graphApiPayload.timeSlotToCheck
	amountOfEaters := graphApiPayload.amountOfEaters

	requestUrl := fmt.Sprintf("https://s-varaukset.fi/api/recommendations/slot/%s/%s/%s/%s", restaurantId, currentDate, timeSlotToCheck, amountOfEaters)
	return requestUrl
}

func (graphApiPayload *raflaamoGraphApiRequestUrl) GenerateGraphApiRequestUrlsForRestaurant(raflaamoTimes *raflaamoTime.RaflaamoTimes) []string {
	requestUrls := make([]string, 0, len(raflaamoTimes.AllFutureGraphApiTimeIntervals))
	for _, graphApiTimeInterval := range raflaamoTimes.AllFutureGraphApiTimeIntervals {
		graphApiPayload.timeSlotToCheck = graphApiTimeInterval
		graphApiRequestUrl := graphApiPayload.getRequestUrlForGraphApi()
		requestUrls = append(requestUrls, graphApiRequestUrl)
	}
	return requestUrls
}
