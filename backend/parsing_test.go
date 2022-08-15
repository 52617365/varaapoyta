package main

import (
	"fmt"
	"strings"
	"testing"
)

// TestGetRestaurantsFromCity | Test to see if JSON parsing works correctly.
func TestGetRestaurantsFromCity(t *testing.T) {
	city := "helsinki"
	restaurants_from_helsinki, err := filter_restaurants_from_city(city)

	if err != nil {
		t.Errorf("Error getting restaurants from city.")
	}
	for _, restaurant := range *restaurants_from_helsinki {
		if !strings.Contains(strings.ToLower(restaurant.Address.Municipality.Fi_FI), strings.ToLower(city)) {
			t.Errorf("restaurant.Address.Municipality.Fi_FI = %s, expected %s", restaurant.Address.Municipality.Fi_FI, "helsinki")
		}
	}
}

// TestGetRestaurantsFromCity | Test to see if JSON parsing works correctly and returns error if nothing found.
func TestGetRestaurantsFromCityThatDoesNotExist(t *testing.T) {
	city := "muumilaakso111"
	restaurants_from_city_that_does_not_exist, err := filter_restaurants_from_city(city)

	if err == nil && len(*restaurants_from_city_that_does_not_exist) > 1 {
		t.Errorf("Expected test to fail but it did not.")
	}
}

// TestBinarySearch | Test to see if binary search algorithm works correctly.
func TestBinarySearch(t *testing.T) {
	times := get_all_possible_reservation_times()
	expected_index := 4
	element_to_find := "0100"
	resulting_index := binary_search(times, element_to_find)
	fmt.Println(resulting_index)

	if expected_index != resulting_index {
		t.Fatalf(`expected index to be %d but it was %d`, expected_index, resulting_index)
	}

}

// todo: make this work.
func TestReturnTimeslotsInbetween(t *testing.T) {
	expected_result_range := [...]string{"0000", "0015", "0030", "0045", "0100"}

	start_time := "2348"
	end_time := "0100"

	time_slots, err := return_time_slots_in_between(start_time, end_time)

	if err != nil {
		t.Fatalf(`TestReturn_time_slots_in_between failed completely with start_time: %s and end_time: %s`, start_time, end_time)
	}
	for index, time_slot := range *time_slots {
		if time_slot != expected_result_range[index] {
			t.Fatalf(`expected time slot to be %s but it was %s`, expected_result_range[index], time_slot)
		}
	}
}

func TestConvert_uneven_minutes_to_even(t *testing.T) {
	test_uneven_number := "1228"
	expected_even_number := "1230"

	even_number := convert_uneven_minutes_to_even(test_uneven_number)

	if even_number != expected_even_number {
		t.Fatalf(`expected even number to be %s but it was %s`, expected_even_number, even_number)
	}

	test_uneven_number2 := "1938"
	expected_even_number2 := "1945"

	even_number2 := convert_uneven_minutes_to_even(test_uneven_number2)

	if even_number2 != expected_even_number2 {
		t.Fatalf(`expected even number to be %s but it was %s`, expected_even_number2, even_number2)
	}
}
