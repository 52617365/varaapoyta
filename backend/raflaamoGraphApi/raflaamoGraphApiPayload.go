/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package raflaamoGraphApi

import (
	"backend/helpers"
	"backend/raflaamoGraphApiTimes"
	"backend/raflaamoRestaurantsApi"
	"fmt"
)

type RequestUrl struct {
	amountOfEaters           string
	timeSlotToCheck          string
	currentDate              string
	IdFromReservationPageUrl string
}

func GetRequestUrl(reservationPageUrl string, amountOfEaters string, currentDate string) *RequestUrl {
	idFromReservationPageUrl := helpers.RegexToMatchRestaurantId.FindString(reservationPageUrl)
	return &RequestUrl{
		amountOfEaters:           amountOfEaters,
		currentDate:              currentDate,
		IdFromReservationPageUrl: idFromReservationPageUrl,
	}
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

func (graphApi *RaflaamoGraphApi) GenerateGraphApiRequestUrlsForRestaurant(restaurant *raflaamoRestaurantsApi.ResponseFields, currentDate string, amountOfEaters string) []string {
	raflaamoGraphApiRequestUrlStruct := GetRequestUrl(restaurant.Links.TableReservationLocalized.FiFi, amountOfEaters, currentDate)
	restaurantsKitchenClosingTime := restaurant.Openingtime.Kitchentime.Ranges[0].End
	graphApiTimeIntervalsFromTheFuture := raflaamoGraphApiTimes.GetAllFutureGraphApiTimeSlots(restaurantsKitchenClosingTime)
	restaurantGraphApiRequestUrls := raflaamoGraphApiRequestUrlStruct.GenerateGraphApiRequestUrlsFromFutureTimeSlots(graphApiTimeIntervalsFromTheFuture)

	return restaurantGraphApiRequestUrls
}
