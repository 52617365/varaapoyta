package main

import (
	"testing"
)

// TestGetRestaurants We expect response len to be over 10.
func TestGetRestaurants(t *testing.T) {
	t.Parallel()
	init_request, _ := init_restaurants()

	response, err := init_request.get()

	if err != nil {
		t.Errorf("unexpected error")
	}

	restaurants_length := len(response)
	if restaurants_length < 10 {
		// Can't check against a static number cuz the amount changes.
		t.Errorf("len(getRestaurants()) = %d, expected %s", restaurants_length, ">10")
	}
}

//func FuzzGetIdFromReservationId(f *testing.F) {
//	f.Add("helsinki", "https://s-varaukset.fi/online/reservation/fi/38?_ga=2.146560948.1092747230.1612503015-489168449.1604043706")
//	f.Fuzz(func(t *testing.T, city string, url string) {
//		placeholder_restaurant := response_fields{
//			Id:          "",
//			Name:        string_field{Fi_FI: ""},
//			Address:     address_fields{Municipality: string_field{Fi_FI: ""}},
//			Features:    features_fields{Accessible: false},
//			Openingtime: opening_fields{Restauranttime: opening_fields_ranges{Ranges: []ranges_times{}}, Kitchentime: opening_fields_ranges{Ranges: []ranges_times{}}},
//			Links:       links_fields{TableReservationLocalized: string_field{Fi_FI: url}, HomepageLocalized: string_field{Fi_FI: ""}},
//		}
//		restaurant_additional_information := additional_information{
//			restaurant: placeholder_restaurant,
//		}
//		id, _ := restaurant_additional_information.get_id_from_reservation_page_url()
//		if !strings.Contains(placeholder_restaurant.Links.TableReservationLocalized.Fi_FI, "https://s-varaukset.fi/online/reservation/fi/") && id != "" {
//			t.Errorf("didnt expect to match regex")
//		}
//	})
//}

// This function is our bottleneck if any.
func TestGetAvailableTables(t *testing.T) {
	t.Parallel()
	amount_of_eaters := 1
	city := "rovaniemi"
	results, err := get_available_tables(city, amount_of_eaters)

	if err != nil {
		t.Errorf("unexpected error")
	}
	if len(results) == 0 {
		t.Errorf("unexpected results length: %d", len(results))
	}
}

func TestGetIdFromReservationPageUrl(t *testing.T) {
	t.Parallel()
	restaurant_url := "https://s-varaukset.fi/online/reservation/fi/38?_ga=2.146560948.1092747230.1612503015-489168449.1604043706"

	expected_id := "38"
	placeholder_restaurant := response_fields{
		Id:          "",
		Name:        &string_field{Fi_FI: ""},
		Address:     &address_fields{Municipality: &string_field{Fi_FI: ""}},
		Openingtime: &opening_fields{Restauranttime: &opening_fields_ranges{Ranges: []ranges_times{}}, Kitchentime: &opening_fields_ranges{Ranges: []ranges_times{}}},
		Links:       &links_fields{TableReservationLocalized: &string_field{Fi_FI: restaurant_url}, HomepageLocalized: &string_field{Fi_FI: ""}},
	}

	restaurant_additional_information := additional_information{
		restaurant: &placeholder_restaurant,
	}

	id, err := restaurant_additional_information.get_id_from_reservation_page_url()

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
		Name:        &string_field{Fi_FI: ""},
		Address:     &address_fields{Municipality: &string_field{Fi_FI: ""}},
		Openingtime: &opening_fields{Restauranttime: &opening_fields_ranges{Ranges: []ranges_times{}}, Kitchentime: &opening_fields_ranges{Ranges: []ranges_times{}}},
		Links:       &links_fields{TableReservationLocalized: &string_field{Fi_FI: restaurant_url}, HomepageLocalized: &string_field{Fi_FI: ""}},
	}
	restaurant_additional_information := additional_information{
		restaurant: &placeholder_restaurant,
	}

	_, err := restaurant_additional_information.get_id_from_reservation_page_url()

	if err == nil {
		t.Errorf("we expected get_id_from_reservation_page_url to throw but it did not.")
	}
}

func BenchmarkGetRestaurants(b *testing.B) {
	for i := 0; i < b.N; i++ {
		init_request, _ := init_restaurants()
		init_request.get()
	}
}

func BenchmarkGetAvailableTables(b *testing.B) {
	for i := 0; i < b.N; i++ {
		amount_of_eaters := 1
		city := "helsinki"
		get_available_tables(city, amount_of_eaters)
	}
}

//func BenchmarkInteractWithApi(b *testing.B) {
//	current_time := get_current_date_and_time()
//	id_from_reservation_page_url := "1769"
//	amount_of_eaters := 1
//	time_slots_to_check_from_graph_api := get_graph_time_slots_from_current_point_forward(current_time.time)
//	for i := 0; i < b.N; i++ {
//		interact_with_api(time_slots_to_check_from_graph_api, id_from_reservation_page_url, current_time.date, amount_of_eaters)
//	}
//}
