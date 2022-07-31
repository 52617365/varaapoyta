package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type structure struct {
	operationName string
	variables     map[string]interface{}
	query         string
}

func getPayload() structure {
	data := structure{
		operationName: "getRestaurantsByLocation",
		variables: map[string]interface{}{
			"first": 470,
			"input": map[string]interface{}{
				"restaurantType": "ALL",
				"locationName":   "Helsinki",
				"feature": map[string]interface{}{
					"rentableVenues": false,
				},
			},
			"after": "eyJmIjowLCJnIjp7ImEiOjYwLjE3MTE2LCJvIjoyNC45MzI1OH19",
		},
		query: "fragment Locales on LocalizedString {fi_FI\n }\n\nfragment Restaurant on Restaurant {\n  id\n  name {\n    ...Locales\n    }\n  urlPath {\n    ...Locales\n     }\n    address {\n    municipality {\n      ...Locales\n       }\n        street {\n      ...Locales\n       }\n       zipCode\n     }\n    features {\n    accessible\n     }\n  openingTime {\n    restaurantTime {\n      ranges {\n        start\n        end\n        endNextDay\n         }\n             }\n    kitchenTime {\n      ranges {\n        start\n        end\n        endNextDay\n              }\n             }\n    }\n  links {\n    tableReservationLocalized {\n      ...Locales\n        }\n    homepageLocalized {\n      ...Locales\n          }\n   }\n     \n}\n\nquery getRestaurantsByLocation($first: Int, $after: String, $input: ListRestaurantsByLocationInput!) {\n  listRestaurantsByLocation(first: $first, after: $after, input: $input) {\n    totalCount\n      edges {\n      ...Restaurant\n        }\n     }\n}",
	}
	return data
}

func getRestaurants() *http.Response {
	data := getPayload()
	dataEncoded, _ := json.Marshal(data)
	resp, err := http.Post("https://api.raflaamo.fi/query", "application/json", bytes.NewBuffer(dataEncoded))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	return resp
}

type response_struct struct {
	id             int         // /data/listRestaurantsByLocation/edges/0/id
	name           string      // /data/listRestaurantsByLocation/edges/0/name/fi_FI
	city           string      // /data/listRestaurantsByLocation/edges/0/address/municipality/fi_FI
	address        string      // /data/listRestaurantsByLocation/edges/0/address/street/fi_FI
	zip_code       string      // /data/listRestaurantsByLocation/edges/0/address/zipCode
	accessible     bool        // /data/listRestaurantsByLocation/edges/0/features/accessible
	opening_hours  interface{} // This will be a tuple containing start and end like (start, end)
	reserving_url  string      // /data/listRestaurantsByLocation/edges/1/links/tableReservationLocalized/fi_FI
	restaurant_url string      // /data/listRestaurantsByLocation/edges/1/links/homepageLocalized/fi_FI
}

func parseRestaurants(city string) {
	restaurants := getRestaurants()

	//TODO: decode the response here.

	fmt.Println(restaurants)
}
