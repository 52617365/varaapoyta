package main

import (
	"strings"
	"testing"
)

// TestGetRestaurants We expect response to be len(470).
func TestGetRestaurants(t *testing.T) {
	t.Parallel()
	restaurants, _ := get_all_restaurants_from_raflaamo_api()

	restaurants_length := len(restaurants)
	if restaurants_length < 10 {
		// Can't check against a static number cuz the amount changes.
		t.Errorf("len(getRestaurants()) = %d, expected %s", restaurants_length, ">10")
	}
}

func FuzzGetIdFromReservationId(f *testing.F) {
	f.Add("helsinki", "https://s-varaukset.fi/online/reservation/fi/38?_ga=2.146560948.1092747230.1612503015-489168449.1604043706")
	f.Fuzz(func(t *testing.T, city string, url string) {
		placeholder_restaurant := response_fields{
			Id:          "",
			Name:        string_field{Fi_FI: ""},
			Urlpath:     string_field{Fi_FI: ""},
			Address:     address_fields{Municipality: string_field{Fi_FI: ""}},
			Features:    features_fields{Accessible: false},
			Openingtime: opening_fields{Restauranttime: opening_fields_ranges{Ranges: []ranges_times{}}, Kitchentime: opening_fields_ranges{Ranges: []ranges_times{}}},
			Links:       links_fields{TableReservationLocalized: string_field{Fi_FI: url}, HomepageLocalized: string_field{Fi_FI: ""}},
		}
		_, err := get_id_from_reservation_page_url(placeholder_restaurant)
		if !strings.Contains(placeholder_restaurant.Links.TableReservationLocalized.Fi_FI, "https://s-varaukset.fi/online/reservation/fi") && err == nil {
			t.Errorf("expected error")
		}
	})
}

// This function is our bottleneck.
func TestGetAvailableTables(t *testing.T) {
	t.Parallel()
	amount_of_eaters := 1
	city := "helsinki"
	restaurants, _ := get_all_restaurants_from_raflaamo_api()
	results := get_available_tables(city, restaurants, amount_of_eaters)

	if len(results) == 0 {
		t.Errorf("unexpected results length: %d", len(results))
	}
}

// reservation_page_url := restaurant.Links.TableReservationLocalized.Fi_FI
func TestGetIdFromReservationPageUrl(t *testing.T) {
	t.Parallel()
	restaurant_url := "https://s-varaukset.fi/online/reservation/fi/38?_ga=2.146560948.1092747230.1612503015-489168449.1604043706"

	expected_id := "38"
	placeholder_restaurant := response_fields{
		Id:          "",
		Name:        string_field{Fi_FI: ""},
		Urlpath:     string_field{Fi_FI: ""},
		Address:     address_fields{Municipality: string_field{Fi_FI: ""}},
		Features:    features_fields{Accessible: false},
		Openingtime: opening_fields{Restauranttime: opening_fields_ranges{Ranges: []ranges_times{}}, Kitchentime: opening_fields_ranges{Ranges: []ranges_times{}}},
		Links:       links_fields{TableReservationLocalized: string_field{Fi_FI: restaurant_url}, HomepageLocalized: string_field{Fi_FI: ""}},
	}

	id, err := get_id_from_reservation_page_url(placeholder_restaurant)

	if err != nil {
		t.Errorf("get_id_from_reservation_page_url threw when we did not expect it to.")
	}

	if id != "38" {
		t.Errorf("get_id_from_reservation_page_url returned %s when we expected %s.", id, expected_id)
	}
}
func TestErrorFromGetIdFromReservationPageUrl(t *testing.T) {
	t.Parallel()
	restaurant_url := "sitethatshouldnotwork.fi"

	placeholder_restaurant := response_fields{
		Id:          "",
		Name:        string_field{Fi_FI: ""},
		Urlpath:     string_field{Fi_FI: ""},
		Address:     address_fields{Municipality: string_field{Fi_FI: ""}},
		Features:    features_fields{Accessible: false},
		Openingtime: opening_fields{Restauranttime: opening_fields_ranges{Ranges: []ranges_times{}}, Kitchentime: opening_fields_ranges{Ranges: []ranges_times{}}},
		Links:       links_fields{TableReservationLocalized: string_field{Fi_FI: restaurant_url}, HomepageLocalized: string_field{Fi_FI: ""}},
	}

	_, err := get_id_from_reservation_page_url(placeholder_restaurant)

	if err == nil {
		t.Errorf("we expected get_id_from_reservation_page_url to throw but it did not.")
	}
}
func BenchmarkGetRestaurants(b *testing.B) {
	for i := 0; i < b.N; i++ {
		get_all_restaurants_from_raflaamo_api()
	}
}
func BenchmarkFilter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		city := "helsinki"
		filter_valid_restaurants_from_city(city)
	}
}
func BenchmarkGetAvailableTables(b *testing.B) {
	for i := 0; i < b.N; i++ {
		amount_of_eaters := 1
		city := "helsinki"
		get_available_tables(city, amount_of_eaters)
	}
}
