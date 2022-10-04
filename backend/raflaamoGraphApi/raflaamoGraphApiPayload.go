package raflaamoGraphApi

import (
	"fmt"
	"regexp"
)

type GraphApiPayload struct {
	reservationPageUrl string
	restaurantId       string
	amountOfEaters     int
	timeSlot           string
	currentDate        string
}

func getRaflaamoGraphApiPayload(reservationPageUrl string, amountOfEaters int, timeSlot string, currentDate string) *GraphApiPayload {
	return &GraphApiPayload{
		reservationPageUrl: reservationPageUrl,
		amountOfEaters:     amountOfEaters,
		timeSlot:           timeSlot,
		currentDate:        currentDate,
	}
}

func (graphApiPayload *GraphApiPayload) getRaflaamoRestaurantIdFromReservationPageUrl() {
	regexToMatchRestaurantId := regexp.MustCompile(`[^fi/]\d+`) // This regex gets the first number match from the TableReservationLocalized JSON field which is the id we want. https://regex101.com/r/NtFMrz/1
	reservationPageUrl := graphApiPayload.reservationPageUrl
	idFromReservationPageUrl := regexToMatchRestaurantId.FindString(reservationPageUrl)
	graphApiPayload.restaurantId = idFromReservationPageUrl
}

func (graphApiPayload *GraphApiPayload) getPayload() string {
	graphApiPayload.getRaflaamoRestaurantIdFromReservationPageUrl()
	restaurantId := graphApiPayload.restaurantId
	currentDate := graphApiPayload.currentDate
	timeSlot := graphApiPayload.timeSlot
	amountOfEaters := graphApiPayload.amountOfEaters

	requestUrl := fmt.Sprintf("https://s-varaukset.fi/api/recommendations/slot/%s/%s/%s/%d", restaurantId, currentDate, timeSlot, amountOfEaters)
	return requestUrl
}
