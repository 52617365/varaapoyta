package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type response_top_level struct {
	Data *response_second_level `json:"data"`
}

type response_second_level struct {
	ListRestaurantsByLocation *response_third_level `json:"listRestaurantsByLocation"`
}

type response_third_level struct {
	Edges *response_fields `json:"edges"`
}

type response_fields struct {
	Id          int              `json:"id"`
	Name        *string_field    `json:"name"`
	Urlpath     *string_field    `json:"urlPath"`
	Address     *address_fields  `json:"address"`
	Features    *features_fields `json:"features"`
	Openingtime *opening_fields  `json:"openingTime"`
	Links       *links_fields    `json:"links"`
}

type string_field struct {
	Fi_FI *string `json:"fi_FI"`
}

type address_fields struct {
	Municipality *string_field `json:"municipality"`
	Street       *string_field `json:"street"`
	Zipcode      *string       `json:"zipCode"`
}

type features_fields struct {
	Accessible bool `json:"accessible"`
}

type opening_fields struct {
	Restauranttime *opening_fields_ranges `json:"restaurantTime"`
	Kitchentime    *opening_fields_ranges `json:"kitchenTime"`
}

type opening_fields_ranges struct {
	Ranges *string `json:"ranges"`
}

type links_fields struct {
	TableReservationLocalized *string_field `json:"tableReservationLocalized"`
	HomepageLocalized         *string_field `json:"homepageLocalized"`
}

func deserialize_response(res **http.Response) *response_top_level {
	response_decoded := &response_top_level{}
	err := json.NewDecoder((*res).Body).Decode(response_decoded)
	if err != nil {
		log.Fatal(err)
	}
	return response_decoded
}
