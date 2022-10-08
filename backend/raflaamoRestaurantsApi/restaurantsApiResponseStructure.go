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
type AvailableTimeSlotsResult struct {
	AvailableTimeSlots chan []string
	Err                error
}

type ResponseFields struct {
	Id              string          `json:"id"`
	Name            *stringField    `json:"name"`
	Address         *addressFields  `json:"address"`
	Openingtime     *openingFields  `json:"openingTime"`
	Links           *linksFields    `json:"links"`
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
	Restauranttime        *openingFieldsRanges `json:"restaurantTime"`
	Kitchentime           *openingFieldsRanges `json:"kitchenTime"`
	kitchenClosingTime    *timeTillKitchenClosingTime
	restaurantClosingTime *timeTillRestaurantClosingTime
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

type timeTillRestaurantClosingTime struct {
	RestaurantClosingHours   int `json:"time_till_restaurant_closed_hours"`
	RestaurantClosingMinutes int `json:"time_till_restaurant_closed_minutes"`
}

// Kitchen closing time - 1 hour determines the time left to reserve.
// This is because the restaurants don't take reservations one hour before the kitchen closes.
type timeTillKitchenClosingTime struct {
	KitchenClosingHours   int `json:"time_left_to_reserve_hours"`
	KitchenClosingMinutes int `json:"time_left_to_reserve_minutes"`
}

// They have the same structure so I'm reusing.
func newKitchenClosingTime(closingHours int, closingMinutes int) *timeTillKitchenClosingTime {
	return &timeTillKitchenClosingTime{KitchenClosingHours: closingHours, KitchenClosingMinutes: closingMinutes}
}
func newRestaurantClosingTime(closingHours int, closingMinutes int) *timeTillRestaurantClosingTime {
	return &timeTillRestaurantClosingTime{RestaurantClosingHours: closingHours, RestaurantClosingMinutes: closingMinutes}
}

func (raflaamoRestaurantsApi *RaflaamoRestaurantsApi) deserializeRaflaamoRestaurantsResponse() (*responseTopLevel, error) {
	responseDecoded := &responseTopLevel{}
	err := json.NewDecoder((raflaamoRestaurantsApi.response).Body).Decode(responseDecoded)
	if err != nil {
		return nil, errors.New("could not deserialize the response body")
	}
	return responseDecoded, nil
}
