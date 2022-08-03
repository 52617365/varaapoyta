package main

import "fmt"

func main() {
	data := getRestaurants()
	for _, value := range data {
		print_string := fmt.Sprintf("%s %s %s %s %s %s %t %s %s", value.Id, *value.Name.Fi_FI, *value.Urlpath.Fi_FI, *value.Address.Municipality.Fi_FI, *value.Address.Street.Fi_FI, *value.Address.Zipcode, value.Features.Accessible, *value.Links.TableReservationLocalized.Fi_FI, *value.Links.HomepageLocalized.Fi_FI)
		fmt.Println(print_string)
	}
}
