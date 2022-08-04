package main

// TODO: Make sure we're not doing any useless copying. (Check return values to see if we return by value).
// TODO: Make endpoints and extract restaurants from a certain country passed in as a parameter.
// (value.Address.Municipality.Fi_FI)
func main() {
	//data := getRestaurants()
	//city := "HELSINKI"
	//restaurants := getRestaurantsFromCity(&*data, &city)
}

// value.Id, *value.Name.Fi_FI, *value.Urlpath.Fi_FI, *value.Address.Municipality.Fi_FI, *value.Address.Street.Fi_FI, *value.Address.Zipcode, value.Features.Accessible, *value.Links.TableReservationLocalized.Fi_FI, *value.Links.HomepageLocalized.Fi_FI
