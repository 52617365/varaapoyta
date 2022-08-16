package main

import (
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
	for _, restaurant := range restaurants_from_helsinki {
		if !strings.Contains(strings.ToLower(restaurant.Address.Municipality.Fi_FI), strings.ToLower(city)) {
			t.Errorf("restaurant.Address.Municipality.Fi_FI = %s, expected %s", restaurant.Address.Municipality.Fi_FI, "helsinki")
		}
	}
}

// TestGetRestaurantsFromCity | Test to see if JSON parsing works correctly and returns error if nothing found.
func TestGetRestaurantsFromCityThatDoesNotExist(t *testing.T) {
	city := "muumilaakso111"
	restaurants_from_city_that_does_not_exist, err := filter_restaurants_from_city(city)

	if err == nil && len(restaurants_from_city_that_does_not_exist) > 1 {
		t.Errorf("Expected test to fail but it did not.")
	}
}

// TestReverseBinarySearch | Test to see if binary search algorithm works correctly.
func TestReverseBinarySearch(t *testing.T) {
	times := get_all_possible_reservation_times()
	expected_index := 91
	element_to_find := "0100"
	resulting_index := reverse_binary_search(times, element_to_find)

	if expected_index != resulting_index {
		t.Errorf(`expected index to be %d but it was %d`, expected_index, resulting_index)
	}

}

func TestReturnTimeslotsInbetween(t *testing.T) {
	expected_result_range := [...]string{"0100", "0045", "0030", "0015", "0000"}

	start_time := "2348"
	end_time := "0100"

	time_slots, err := return_time_slots_in_between(start_time, end_time)

	if err != nil {
		t.Errorf(`TestReturn_time_slots_in_between failed completely with start_time: %s and end_time: %s`, start_time, end_time)
	}

	for index, _ := range time_slots {
		if time_slots[index] != expected_result_range[index] {
			t.Errorf(`expected time slot to be %s but it was %s`, time_slots[index], expected_result_range[index])
		}
	}
}

func TestReturnTimeslotsInbetween2(t *testing.T) {
	expected_result_range := [...]string{"1800",
		"1745", "1730", "1715", "1700",
		"1645", "1630", "1615", "1600",
		"1545", "1530", "1515", "1500",
		"1445", "1430", "1415", "1400",
		"1345", "1330", "1315", "1300",
		"1245", "1230", "1215", "1200",
		"1145", "1130", "1115", "1100",
		"1045", "1030", "1015", "1000",
		"0945", "0930", "0915", "0900",
		"0845", "0830", "0815", "0800",
		"0745", "0730", "0715", "0700",
		"0645", "0630", "0615", "0600",
		"0545", "0530", "0515", "0500",
		"0445", "0430", "0415", "0400",
		"0345", "0330", "0315", "0300",
		"0245", "0230", "0215", "0200",
		"0145", "0130", "0115", "0100",
		"0045", "0030", "0015", "0000"}

	start_time := "1800"
	end_time := "2359"

	// FIX: it gets stuck in this function.
	time_slots, err := return_time_slots_in_between(start_time, end_time)

	if err != nil {
		t.Errorf(`TestReturn_time_slots_in_between failed completely with start_time: %s and end_time: %s`, start_time, end_time)
	}

	for index, _ := range time_slots {
		if time_slots[index] != expected_result_range[index] {
			t.Errorf(`expected time slot to be %s but it was %s`, time_slots[0], expected_result_range[0])
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
