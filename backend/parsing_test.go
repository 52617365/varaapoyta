package main

import (
	"strings"
	"testing"
)

// TestGetRestaurantsFromCity | Test to see if JSON parsing works correctly.
func TestGetRestaurantsFromCity(t *testing.T) {
	t.Parallel()
	city := "helsinki"
	restaurants_from_helsinki, err := filter_valid_restaurants_from_city(city)

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
	restaurants_from_city_that_does_not_exist, err := filter_valid_restaurants_from_city(city)

	if err == nil && len(restaurants_from_city_that_does_not_exist) > 1 {
		t.Errorf("Expected test to fail but it did not.")
	}
}
