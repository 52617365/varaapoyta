/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package raflaamoGraphApi

import (
	"backend/helpers"
	"backend/raflaamoGraphApiTimes"
	"errors"
	"fmt"
)

type RequestUrl struct {
	amountOfEaters           string
	timeSlotToCheck          string
	currentDate              string
	IdFromReservationPageUrl string
}

func GetRequestUrl(reservationPageUrl string, amountOfEaters string, currentDate string) (*RequestUrl, error) {
	regexMatchGroups := helpers.RegexToMatchRestaurantId.FindAllStringSubmatch(reservationPageUrl, -1)
	idFromReservationPageUrl := regexMatchGroups[0][1]
	if idFromReservationPageUrl == "" {
		return nil, fmt.Errorf("[GetRequestUrl] - %w", errors.New("there was an error matching regex"))
	}
	return &RequestUrl{
		amountOfEaters:           amountOfEaters,
		currentDate:              currentDate,
		IdFromReservationPageUrl: idFromReservationPageUrl,
	}, nil
}

func (graphApiPayload *RequestUrl) getRequestUrlForGraphApi(timeSlotToCheck string) string {
	restaurantId := graphApiPayload.IdFromReservationPageUrl
	currentDate := graphApiPayload.currentDate
	amountOfEaters := graphApiPayload.amountOfEaters

	requestUrl := fmt.Sprintf("https://s-varaukset.fi/api/recommendations/slot/%s/%s/%s/%s", restaurantId, currentDate, timeSlotToCheck, amountOfEaters)
	return requestUrl
}

func (graphApiPayload *RequestUrl) GenerateGraphApiRequestUrlsFromFutureTimeSlots(graphApiTimeSlotsFromTheFuture []string) []string {
	requestUrls := make([]string, 0, len(graphApiTimeSlotsFromTheFuture))
	for _, graphApiTimeInterval := range graphApiTimeSlotsFromTheFuture {
		graphApiRequestUrl := graphApiPayload.getRequestUrlForGraphApi(graphApiTimeInterval)
		requestUrls = append(requestUrls, graphApiRequestUrl)
	}
	return requestUrls
}

func (graphApi *RaflaamoGraphApi) GenerateGraphApiRequestUrlsForRestaurant(currentTime int64, raflaamoGraphApiRequestUrlStruct *RequestUrl) []string {
	graphApiTimeIntervalsFromTheFuture := raflaamoGraphApiTimes.GetAllFutureGraphApiTimeSlots(currentTime)
	restaurantGraphApiRequestUrls := raflaamoGraphApiRequestUrlStruct.GenerateGraphApiRequestUrlsFromFutureTimeSlots(graphApiTimeIntervalsFromTheFuture)

	return restaurantGraphApiRequestUrls
}
