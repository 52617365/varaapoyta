package main

import (
	"strings"
	"testing"
)

//func Fuzz_convert_uneven_minutes_to_even(f *testing.F) {
//	even_minutes := []string{"15", "30", "45", "00"}
//	f.Add("1500")
//	f.Add("0039")
//	f.Add("A039")
//	f.Fuzz(func(t *testing.T, time_to_convert string) {
//		even_number, convert_err := convert_uneven_time_to_even(time_to_convert)
//
//		// if it can't be parsed into an integer and we did not get an error.
//		if _, err := strconv.ParseInt(time_to_convert, 10, 64); err != nil {
//			if convert_err == nil {
//				t.Errorf(`expected an error with number: %s`, time_to_convert)
//			}
//		}
//
//		if _, err := strconv.ParseInt(time_to_convert, 10, 64); err == nil {
//			if convert_err != nil {
//				t.Errorf(`did not expect an error with time: %s, error is: `, time_to_convert, convert_err)
//			}
//			if even_number == "" {
//				t.Errorf(`expected a result with number: %s`, time_to_convert)
//			}
//		}
//		result_minutes := even_number[len(even_number)-2:]
//		if !slices.Contains(even_minutes, result_minutes) {
//			t.Errorf(`Expected minutes to be 15, 30, 45 or 00 but it was %s`, result_minutes)
//		}
//	})
//}

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
