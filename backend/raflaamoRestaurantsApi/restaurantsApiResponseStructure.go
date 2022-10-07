package raflaamoRestaurantsApi

import (
	"encoding/json"
	"errors"
)

type responseTopLevel struct {
	Data *responseSecondLevel `json:"data"`
}

type responseSecondLevel struct {
	ListRestaurantsByLocation *responseThirdLevel `json:"listRestaurantsByLocation"`
}

type responseThirdLevel struct {
	Edges []ResponseFields `json:"edges"`
}
type AvailableTimeSlotsResult struct {
	AvailableTimeSlots chan []string
	Err                error
}

type ResponseFields struct {
	Id          string         `json:"id"`
	Name        *stringField   `json:"name"`
	Address     *addressFields `json:"address"`
	Openingtime *openingFields `json:"openingTime"`
	Links       *linksFields   `json:"links"`
	//AvailableTimeSlotsBuffer *AvailableTimeSlotsResult `json:"available_time_slots"` // This will be populated later on when we iterate this list and get all raflaamoTime slots.
	//AvailableTimeSlotsBuffer []string `json:"available_time_slots"` // This will be populated later on when we iterate this list and get all raflaamoTime slots.
	GraphApiResults *GraphApiResult `json:"available_time_slots"` // This will be populated later on when we iterate this list and get all raflaamoTime slots.
}

type stringField struct {
	FiFi string `json:"fi_FI"`
}

type addressFields struct {
	Municipality *stringField `json:"municipality"`
	Street       *stringField `json:"street"`
	Zipcode      string       `json:"zipCode"`
}

type openingFields struct {
	Restauranttime                  *openingFieldsRanges `json:"restaurantTime"`
	Kitchentime                     *openingFieldsRanges `json:"kitchenTime"`
	TimeTillRestaurantClosedHours   int                  `json:"time_till_restaurant_closed_hours"`
	TimeTillRestaurantClosedMinutes int                  `json:"time_till_restaurant_closed_minutes"`
	TimeLeftToReserveHours          int                  `json:"time_left_to_reserve_hours"`
	TimeLeftToReserveMinutes        int                  `json:"time_left_to_reserve_minutes"`
}

type openingFieldsRanges struct {
	Ranges []rangesTimes `json:"ranges"`
}

type rangesTimes struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type linksFields struct {
	TableReservationLocalized   *stringField `json:"tableReservationLocalized"`
	TableReservationLocalizedId string       `json:"tableReservationLocalizedId"`
	HomepageLocalized           *stringField `json:"homepageLocalized"`
}

type GraphApiResult struct {
	AvailableTimeSlotsBuffer chan string
	Err                      chan error
}

func (raflaamoRestaurantsApi *RaflaamoRestaurantsApi) deserializeRaflaamoRestaurantsResponse() (*responseTopLevel, error) {
	responseDecoded := &responseTopLevel{}
	err := json.NewDecoder((raflaamoRestaurantsApi.response).Body).Decode(responseDecoded)
	if err != nil {
		return nil, errors.New("could not deserialize the response body")
	}
	return responseDecoded, nil
}
