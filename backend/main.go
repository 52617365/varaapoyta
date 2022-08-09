package main

import "fmt"

// TODO: Figure out what tables are free to reserve. (Format string correctly and make request and save result.)
// TODO: Make endpoints and extract restaurants from a certain country passed in as a parameter.

// (value.Address.Municipality.Fi_FI)
func main() {
	city := "Rovaniemi"
	restaurants, err := getRestaurantsFromCity(&city)
	if err != nil {
		fmt.Println("Could not find any restaurants.")
		return
	}
	available_ones := getAvailableTables(restaurants)
	for _, available := range *available_ones {
		fmt.Println(available)
	}
	// city := "muumilaakso"
	// restaurants, err := getRestaurantsFromCity(&city)
	// if err != nil {
	// 	fmt.Println("Could not find any restaurants.")
	// 	return
	// }
	// for _, v := range *restaurants {
	// 	fmt.Println(*v.Id)
	// }

}

// *value.Id, *value.Name.Fi_FI, *value.Urlpath.Fi_FI, *value.Address.Municipality.Fi_FI, *value.Address.Street.Fi_FI, *value.Address.Zipcode, value.Features.Accessible, *value.Links.TableReservationLocalized.Fi_FI, *value.Links.HomepageLocalized.Fi_FI
