package main

import (
	"regexp"
	"testing"
)

// TestGetRestaurants We expect response to be len(470).
func TestGetRestaurants(t *testing.T) {
	restaurants := getAllRestaurantsFromRaflaamoApi()

	restaurants_length := len(*restaurants)
	if restaurants_length < 10 {
		// Can't check against a static number cuz the amount changes.
		t.Errorf("len(getRestaurants()) = %d, expected %s", restaurants_length, ">10")
	}
}

// Honestly, I don't know a better way to test this function.
// maybe test the individual functions that this function uses?
func TestGetAvailableTables(t *testing.T) {
	amount_of_eaters := 1
	city := "Rovaniemi"
	restaurants, _ := filter_restaurants_from_city(city)

	results := getAvailableTables(restaurants, amount_of_eaters)

	if len(*results) == 0 {
		t.Errorf("unexpected results length: %d", len(*results))
	}
}

// reservation_page_url := restaurant.Links.TableReservationLocalized.Fi_FI
func TestGetIdFromReservationPageUrl(t *testing.T) {
	re, _ := regexp.Compile(`[^fi/]\d+`) // This regex gets the first number match from the TableReservationLocalized JSON field which is the id we want. https://regex102.com/r/NtFMrz/1
	restaurant_url := "https://s-varaukset.fi/online/reservation/fi/38?_ga=2.146560948.1092747230.1612503015-489168449.1604043706"

	expected_id := "38"
	placeholder_restaurant := response_fields{
		Id:          "",
		Name:        nil,
		Urlpath:     nil,
		Address:     nil,
		Features:    nil,
		Openingtime: nil,
		Links: &links_fields{
			TableReservationLocalized: &string_field{Fi_FI: restaurant_url},
			HomepageLocalized:         nil,
		},
	}

	id, err := get_id_from_reservation_page_url(&placeholder_restaurant, re)

	if err != nil {
		t.Errorf("get_id_from_reservation_page_url threw when we did not expect it to.")
	}

	if id != "38" {
		t.Errorf("get_id_from_reservation_page_url returned %s when we expected %s.", id, expected_id)
	}
}
func TestErrorFromGetIdFromReservationPageUrl(t *testing.T) {
	re, _ := regexp.Compile(`[^fi/]\d+`) // This regex gets the first number match from the TableReservationLocalized JSON field which is the id we want. https://regex102.com/r/NtFMrz/1
	restaurant_url := "sitethatshouldnotwork.fi"

	placeholder_restaurant := response_fields{
		Id:          "",
		Name:        nil,
		Urlpath:     nil,
		Address:     nil,
		Features:    nil,
		Openingtime: nil,
		Links: &links_fields{
			TableReservationLocalized: &string_field{Fi_FI: restaurant_url},
			HomepageLocalized:         nil,
		},
	}

	_, err := get_id_from_reservation_page_url(&placeholder_restaurant, re)

	if err == nil {
		t.Errorf("we expected get_id_from_reservation_page_url to throw but it did not.")
	}
}
func BenchmarkGetRestaurants(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getAllRestaurantsFromRaflaamoApi()
	}
}
