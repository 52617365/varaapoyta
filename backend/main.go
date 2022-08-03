package main

import (
	"fmt"
	"strings"
)

// TODO: Make endpoints and extract restaurants from a certain country passed in as a parameter.
// (value.Address.Municipality.Fi_FI)
func main() {
	data := getRestaurants()
	city := "HELSINKI"
	restaurants := getRestaurantsFromCity(&data, &city)
	for _, restaurant := range restaurants {
		print_string := fmt.Sprintf("%s %s", *restaurant.Name.Fi_FI, *restaurant.Address.Municipality.Fi_FI)
		fmt.Println(print_string)
	}
	//	for _, value := range data {
	//		print_string := fmt.Sprintf("%s %s %s %s %s %s %t %s %s", value.Id, *value.Name.Fi_FI, *value.Urlpath.Fi_FI, *value.Address.Municipality.Fi_FI, *value.Address.Street.Fi_FI, *value.Address.Zipcode, value.Features.Accessible, *value.Links.TableReservationLocalized.Fi_FI, *value.Links.HomepageLocalized.Fi_FI)
	//		fmt.Println(print_string)
	//	}
}

// TODO: capture result.
func getRestaurantsFromCity(restaurants *[]response_fields, city *string) []response_fields {
	captured_restaurants := []response_fields{}
	for _, restaurant := range *restaurants {
		if strings.Contains(strings.ToLower(*restaurant.Address.Municipality.Fi_FI), strings.ToLower(*city)) {
			captured_restaurants = append(captured_restaurants, restaurant)
		}
	}
	return captured_restaurants
}
