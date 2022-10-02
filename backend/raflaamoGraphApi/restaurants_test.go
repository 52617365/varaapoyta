package raflaamoGraphApi

// import (
// 	"testing"
// )

// // TestGetRestaurants We expect response len to be over 10.
// func TestGetRestaurants(t *testing.T) {
// 	t.Parallel()
// 	initRequest, _ := init_restaurants()

// 	response, err := initRequest.get()

// 	if err != nil {
// 		t.Errorf("unexpected error")
// 	}

// 	restaurantsLength := len(response)
// 	if restaurantsLength < 10 {
// 		// Can't check against a static number cuz the amount changes.
// 		t.Errorf("len(getRestaurants()) = %d, expected %s", restaurantsLength, ">10")
// 	}
// }

// // This function is our bottleneck if any.
// func TestGetAvailableTables(t *testing.T) {
// 	t.Parallel()
// 	amountOfEaters := 1
// 	city := "rovaniemi"
// 	results, err := getAvailableTables(city, amountOfEaters)

// 	if err != nil {
// 		t.Errorf("unexpected error")
// 	}
// 	if len(results) == 0 {
// 		t.Errorf("unexpected results length: %d", len(results))
// 	}
// }

// func TestGetIdFromReservationPageUrl(t *testing.T) {
// 	t.Parallel()
// 	restaurantUrl := "https://s-varaukset.fi/online/reservation/fi/38?_ga=2.146560948.1092747230.1612503015-489168449.1604043706"

// 	expectedId := "38"
// 	placeholderRestaurant := response_fields{
// 		Id:          "",
// 		Name:        &string_field{Fi_FI: ""},
// 		Address:     &address_fields{Municipality: &string_field{Fi_FI: ""}},
// 		Openingtime: &opening_fields{Restauranttime: &opening_fields_ranges{Ranges: []ranges_times{}}, Kitchentime: &opening_fields_ranges{Ranges: []ranges_times{}}},
// 		Links:       &links_fields{TableReservationLocalized: &string_field{Fi_FI: restaurantUrl}, HomepageLocalized: &string_field{Fi_FI: ""}},
// 	}

// 	restaurantAdditionalInformation := AdditionalInformationAboutRestaurant{
// 		restaurant: &placeholderRestaurant,
// 	}

// 	id, err := restaurantAdditionalInformation.get_id_from_reservation_page_url()

// 	if err != nil {
// 		t.Errorf("get_id_from_reservation_page_url threw when we did not expect it to.")
// 	}

// 	if id != "38" {
// 		t.Errorf("get_id_from_reservation_page_url returned %s when we expected %s.", id, expectedId)
// 	}
// }
// func TestErrorFromGetIdFromReservationPageUrl(t *testing.T) {
// 	t.Parallel()
// 	restaurantUrl := "sitethatshouldnotwork.fi"

// 	placeholderRestaurant := response_fields{
// 		Id:          "",
// 		Name:        &string_field{Fi_FI: ""},
// 		Address:     &address_fields{Municipality: &string_field{Fi_FI: ""}},
// 		Openingtime: &opening_fields{Restauranttime: &opening_fields_ranges{Ranges: []ranges_times{}}, Kitchentime: &opening_fields_ranges{Ranges: []ranges_times{}}},
// 		Links:       &links_fields{TableReservationLocalized: &string_field{Fi_FI: restaurantUrl}, HomepageLocalized: &string_field{Fi_FI: ""}},
// 	}
// 	restaurantAdditionalInformation := AdditionalInformationAboutRestaurant{
// 		restaurant: &placeholderRestaurant,
// 	}

// 	_, err := restaurantAdditionalInformation.get_id_from_reservation_page_url()

// 	if err == nil {
// 		t.Errorf("we expected get_id_from_reservation_page_url to throw but it did not.")
// 	}
// }

// func BenchmarkGetRestaurants(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		initRequest, _ := init_restaurants()
// 		initRequest.get()
// 	}
// }

// func BenchmarkGetAvailableTables(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		amount_of_eaters := 1
// 		city := "helsinki"
// 		getAvailableTables(city, amount_of_eaters)
// 	}
// }

// //func BenchmarkInteractWithApi(b *testing.B) {
// //	current_time := get_current_date_and_time()
// //	id_from_reservation_page_url := "1769"
// //	amount_of_eaters := 1
// //	time_slots_to_check_from_graph_api := get_graph_time_slots_from_current_point_forward(current_time.timeUtils)
// //	for i := 0; i < b.N; i++ {
// //		interact_with_api(time_slots_to_check_from_graph_api, id_from_reservation_page_url, current_time.date, amount_of_eaters)
// //	}
// //}
