package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type payload_structure struct {
	Operationname string              `json:"operationName"`
	Variables     variables_structure `json:"variables"`
	Query         string              `json:"query"`
}
type variables_structure struct {
	First int             `json:"first"`
	Input input_structure `json:"input"`
	After string          `json:"after"`
}
type input_structure struct {
	RestaurantType string            `json:"restaurantType"`
	Locationname   string            `json:"locationName"`
	Feature        feature_structure `json:"feature"`
}
type feature_structure struct {
	Rentablevenues bool `json:"rentableVenues"`
}

func generate_and_serialize_payload() []byte {
	// Making the post request payload here to raflaamo api.
	payload := &payload_structure{
		Operationname: "getRestaurantsByLocation",
		Variables: variables_structure{
			First: 470,
			Input: input_structure{
				RestaurantType: "ALL",
				Locationname:   "Helsinki",
				Feature:        feature_structure{Rentablevenues: false},
			},
			After: "eyJmIjowLCJnIjp7ImEiOjYwLjE3MTE2LC",
		},
		Query: `fragment Locales on LocalizedString {fi_FI\n }\n\nfragment Restaurant on Restaurant {\n  id\n  name {\n    ...Locales\n    }\n  urlPath {\n    ...Locales\n     }\n    address {\n    municipality {\n      ...Locales\n       }\n        street {\n      ...Locales\n       }\n       zipCode\n     }\n    features {\n    accessible\n     }\n  openingTime {\n    restaurantTime {\n      ranges {\n        start\n        end\n        endNextDay\n         }\n             }\n    kitchenTime {\n      ranges {\n        start\n        end\n        endNextDay\n              }\n             }\n    }\n  links {\n    tableReservationLocalized {\n      ...Locales\n        }\n    homepageLocalized {\n      ...Locales\n          }\n   }\n     \n}\n\nquery getRestaurantsByLocation($first: Int, $after: String, $input: ListRestaurantsByLocationInput!) {\n  listRestaurantsByLocation(first: $first, after: $after, input: $input) {\n    totalCount\n      edges {\n      ...Restaurant\n        }\n     }\n}`,
	}

	// Turning map into json string thats ready to send here.
	json_bytes, err := json.Marshal(payload)
	fmt.Println(string(json_bytes))

	if err != nil {
		log.Fatal(err)
	}

	return json_bytes
}
