package raflaamoGraphApi

import "fmt"

type GraphApiPayload struct {
	restaurantId   string
	amountOfEaters int
	timeSlot       string
	currentDate    string
}

func getRaflaamoGraphApiPayload(restaurantId string, amountOfEaters int, timeSlot string, currentDate string) *GraphApiPayload {
	return &GraphApiPayload{
		restaurantId:   restaurantId,
		amountOfEaters: amountOfEaters,
		timeSlot:       timeSlot,
		currentDate:    currentDate,
	}
}

func (graphApiPayload *GraphApiPayload) getPayload() string {
	restaurantId := graphApiPayload.restaurantId
	amountOfEaters := graphApiPayload.amountOfEaters
	timeSlot := graphApiPayload.timeSlot
	currentDate := graphApiPayload.currentDate

	requestUrl := fmt.Sprintf("https://s-varaukset.fi/api/recommendations/slot/%s/%s/%s/%d", restaurantId, currentDate, timeSlot, amountOfEaters)
	return requestUrl
}
