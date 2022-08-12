package main

import "fmt"

// FIX: Make checking of available tables work.
// TODO: Make endpoints and extract restaurants from a certain country passed in as a parameter.

func main() {
	city := "Helsinki"
	restaurants, err := getRestaurantsFromCity(&city)
	if err != nil {
		fmt.Println("Could not find any restaurants.")
		return
	}
	getAvailableTables(restaurants)
	fmt.Println(restaurants)
}

// *value.Id, *value.Name.Fi_FI, *value.Urlpath.Fi_FI, *value.Address.Municipality.Fi_FI, *value.Address.Street.Fi_FI, *value.Address.Zipcode, value.Features.Accessible, *value.Links.TableReservationLocalized.Fi_FI, *value.Links.HomepageLocalized.Fi_FI
