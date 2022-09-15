package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type response_top_level struct {
	Data response_second_level `json:"data"`
}

type response_second_level struct {
	ListRestaurantsByLocation response_third_level `json:"listRestaurantsByLocation"`
}

type response_third_level struct {
	Edges []response_fields `json:"edges"`
}

type response_fields struct {
	Id                   string          `json:"id"`
	Name                 string_field    `json:"name"`
	Address              address_fields  `json:"address"`
	Features             features_fields `json:"features"`
	Openingtime          opening_fields  `json:"openingTime"`
	Links                links_fields    `json:"links"`
	Available_time_slots []string        `json:"available_time_slots"` // This will be populated later on when we iterate this list and get all time slots.
}

type string_field struct {
	Fi_FI string `json:"fi_FI"`
}

type address_fields struct {
	Municipality string_field `json:"municipality"`
	Street       string_field `json:"street"`
	Zipcode      string       `json:"zipCode"`
}

type features_fields struct {
	Accessible bool `json:"accessible"`
}

type opening_fields struct {
	Restauranttime                      opening_fields_ranges `json:"restaurantTime"`
	Kitchentime                         opening_fields_ranges `json:"kitchenTime"`
	Time_till_restaurant_closed_hours   int                   `json:"time_till_restaurant_closed_hours"`
	Time_till_restaurant_closed_minutes int                   `json:"time_till_restaurant_closed_minutes"`
	Time_left_to_reserve_hours          int                   `json:"time_left_to_reserve_hours"`
	Time_left_to_reserve_minutes        int                   `json:"time_left_to_reserve_minutes"`
}

type opening_fields_ranges struct {
	// Ranges interface{} `json:"ranges"`
	Ranges []ranges_times `json:"ranges"`
}

type ranges_times struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type links_fields struct {
	TableReservationLocalized   string_field `json:"tableReservationLocalized"`
	TableReservationLocalizedId string       `json:"tableReservationLocalizedId"`
	HomepageLocalized           string_field `json:"homepageLocalized"`
}

// Tries to deserialize the response from the raflaamo API and returns an error if it fails.
func deserialize_response(res **http.Response) (response_top_level, error) {
	response_decoded := &response_top_level{}
	err := json.NewDecoder((*res).Body).Decode(response_decoded)
	if err != nil {
		return *response_decoded, errors.New("could not deserialize the response body")
	}
	return *response_decoded, nil
}
