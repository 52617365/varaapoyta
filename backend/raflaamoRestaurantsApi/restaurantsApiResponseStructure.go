/*
 * Copyright (c) 2022. Rasmus Mäki
 */

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

type ResponseFields struct {
	Id                 string          `json:"id"`
	Name               *StringField    `json:"name"`
	Address            *AddressFields  `json:"address"`
	Openingtime        *OpeningFields  `json:"openingTime"`
	Links              *linksFields    `json:"links"`
	GraphApiResults    *GraphApiResult `json:"-"` // This will be populated later on when we iterate this list and get all raflaamoTimes slots.
	AvailableTimeSlots []string        `json:"available_time_slots"`
}

type StringField struct {
	FiFi string `json:"fi_FI"`
}

type AddressFields struct {
	Municipality *StringField `json:"municipality"`
	Street       *StringField `json:"street"`
	Zipcode      string       `json:"zipCode"`
}

type OpeningFields struct {
	Restauranttime                  *openingFieldsRanges `json:"restaurantTime"`
	Kitchentime                     *openingFieldsRanges `json:"kitchenTime"`
	TimeTillRestaurantClosedHours   int                  `json:"time_till_restaurant_closed_hours"`
	TimeTillRestaurantClosedMinutes int                  `json:"time_till_restaurant_closed_minutes"`
	TimeLeftToReserveHours          int                  `json:"time_left_to_reserve_hours"`
	TimeLeftToReserveMinutes        int                  `json:"time_left_to_reserve_minutes"`
}

type openingFieldsRanges struct {
	Ranges []RangesTimes `json:"ranges"`
}

type RangesTimes struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type linksFields struct {
	TableReservationLocalized   *StringField `json:"tableReservationLocalized"`
	TableReservationLocalizedId string       `json:"tableReservationLocalizedId"`
	HomepageLocalized           *StringField `json:"homepageLocalized"`
}

type GraphApiResult struct {
	AvailableTimeSlotsBuffer chan string `json:"-"`
	Err                      chan error  `json:"-"`
}

func (raflaamoRestaurantsApi *RaflaamoRestaurantsApi) deserializeRaflaamoRestaurantsResponse() (*responseTopLevel, error) {
	responseDecoded := &responseTopLevel{}
	err := json.NewDecoder((raflaamoRestaurantsApi.response).Body).Decode(responseDecoded)
	if err != nil {
		return nil, errors.New("could not deserialize the response body")
	}
	return responseDecoded, nil
}
