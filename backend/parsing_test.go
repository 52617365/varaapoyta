package main

import (
	"strings"
	"testing"
)

// TestGetRestaurantsFromCity | Test to see if JSON parsing works correctly.
func TestGetRestaurantsFromCity(t *testing.T) {
	t.Parallel()
	city := "helsinki"
	restaurants_from_helsinki, err := filter_restaurants_from_city(city)

	if err != nil {
		t.Errorf("Error getting restaurants from city.")
	}
	for _, restaurant := range restaurants_from_helsinki {
		if !strings.Contains(strings.ToLower(restaurant.Address.Municipality.Fi_FI), strings.ToLower(city)) {
			t.Errorf("restaurant.Address.Municipality.Fi_FI = %s, expected %s", restaurant.Address.Municipality.Fi_FI, "helsinki")
		}
	}
}

// TestGetRestaurantsFromCity | Test to see if JSON parsing works correctly and returns error if nothing found.
func TestGetRestaurantsFromCityThatDoesNotExist(t *testing.T) {
	t.Parallel()
	city := "muumilaakso111"
	restaurants_from_city_that_does_not_exist, err := filter_restaurants_from_city(city)

	if err == nil && len(restaurants_from_city_that_does_not_exist) > 1 {
		t.Errorf("Expected test to fail but it did not.")
	}
}

// // TestReverseBinarySearch | Test to see if binary search algorithm works correctly.
// func TestBinarySearch(t *testing.T) {
// 	t.Parallel()
// 	times := get_all_reservation_times("0200") // in reality, it's not all because we need to consider restaurants closing time.
// 	expected_index := 4
// 	element_to_find := "0100"
// 	resulting_index := binary_search(times, element_to_find)

// 	if expected_index != resulting_index {
// 		t.Errorf(`expected index to be %d but it was %d`, expected_index, resulting_index)
// 	}
// }
