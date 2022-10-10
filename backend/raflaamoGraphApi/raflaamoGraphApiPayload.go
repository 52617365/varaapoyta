/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package raflaamoGraphApi

import (
	"backend/raflaamoTime"
	"backend/regex"
	"fmt"
)

type RequestUrl struct {
	amountOfEaters           string
	timeSlotToCheck          string
	currentDate              string
	IdFromReservationPageUrl string
}

func GetRaflaamoGraphApiRequestUrl(reservationPageUrl string, amountOfEaters string, currentDate string) *RequestUrl {
	idFromReservationPageUrl := regex.RegexToMatchRestaurantId.FindString(reservationPageUrl)
	return &RequestUrl{
		amountOfEaters:           amountOfEaters,
		currentDate:              currentDate,
		IdFromReservationPageUrl: idFromReservationPageUrl,
	}
}

func (graphApiPayload *RequestUrl) getRequestUrlForGraphApi() string {
	restaurantId := graphApiPayload.IdFromReservationPageUrl
	currentDate := graphApiPayload.currentDate
	timeSlotToCheck := graphApiPayload.timeSlotToCheck
	amountOfEaters := graphApiPayload.amountOfEaters

	requestUrl := fmt.Sprintf("https://s-varaukset.fi/api/recommendations/slot/%s/%s/%s/%s", restaurantId, currentDate, timeSlotToCheck, amountOfEaters)
	return requestUrl
}

func (graphApiPayload *RequestUrl) GenerateGraphApiRequestUrlsForRestaurant(raflaamoTimes *raflaamoTime.RaflaamoTimes) []string {
	requestUrls := make([]string, 0, len(raflaamoTimes.AllFutureGraphApiTimeIntervals))
	for _, graphApiTimeInterval := range raflaamoTimes.AllFutureGraphApiTimeIntervals {
		graphApiPayload.timeSlotToCheck = graphApiTimeInterval
		graphApiRequestUrl := graphApiPayload.getRequestUrlForGraphApi()
		requestUrls = append(requestUrls, graphApiRequestUrl)
	}
	return requestUrls
}
