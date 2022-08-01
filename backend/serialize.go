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
		Query: `fragment Locales on LocalizedString {fi_FI} fragment Restaurant on Restaurant {  id  name {...Locales}  urlPath {...Locales}    address {    municipality {      ...Locales       }        street {      ...Locales       }       zipCode     }    features {    accessible     }  openingTime {    restaurantTime {      ranges {        start        end        endNextDay         }             }    kitchenTime {      ranges {        start        end        endNextDay              }             }    }  links {    tableReservationLocalized {      ...Locales        }    homepageLocalized {      ...Locales          }   }     }query getRestaurantsByLocation($first: Int, $after: String, $input: ListRestaurantsByLocationInput!){  listRestaurantsByLocation(first: $first, after: $after, input: $input) {    totalCount      edges {      ...Restaurant        }     }}`,
	}

	json_bytes, err := json.Marshal(payload)
	fmt.Println(string(json_bytes))

	if err != nil {
		log.Fatal(err)
	}

	return json_bytes
}
