package main

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

// TestGetRestaurantsFromCity | Test to see if JSON parsing works correctly.
func TestGetRestaurantsFromCity(t *testing.T) {
	city := "helsinki"
	restaurants_from_helsinki, err := filter_restaurants_from_city(&city)

	if err != nil {
		t.Errorf("Error getting restaurants from city.")
	}
	for _, restaurant := range *restaurants_from_helsinki {
		if !strings.Contains(strings.ToLower(*restaurant.Address.Municipality.Fi_FI), strings.ToLower(city)) {
			t.Errorf("restaurant.Address.Municipality.Fi_FI = %s, expected %s", *restaurant.Address.Municipality.Fi_FI, "helsinki")
		}
	}
}

// TestGetRestaurantsFromCity | Test to see if JSON parsing works correctly and returns error if nothing found.
func TestGetRestaurantsFromCityThatDoesNotExist(t *testing.T) {
	city := "muumilaakso111"
	restaurants_from_city_that_does_not_exist, err := filter_restaurants_from_city(&city)

	if err == nil && len(*restaurants_from_city_that_does_not_exist) > 1 {
		t.Errorf("Expected test to fail but it did not.")
	}
}

// TestGetAllPossibleTimes | Test to see if function returns all possible times forward from current time.
func TestGetAllPossibleTimes(t *testing.T) {
	currentTime, _ := strconv.Atoi(strings.ReplaceAll(getCurrentTime(), ":", ""))
	allPossibleTimes := getAllPossibleTimes()

	for _, time := range *allPossibleTimes {
		time_int, _ := strconv.Atoi(strings.ReplaceAll(time, ":", ""))
		if currentTime > time_int {
			t.Fatalf(`allPossibleTimes() > time_int returned %d, expected more than %d.`, time_int, currentTime)
		}
	}
}

// TestFormatTimesToString | Test to see if function converts times stored in ints to correctly formatted string times.
func TestFormatTimesToString(t *testing.T) {
	times_int := []int{900, 1000}

	times_strings := convert_times_to_string(times_int)

	if times_strings[0] != "09:00" {
		t.Fatalf(`formatTimesToString() converted 900 to %s, expected %s.`, times_strings[0], "09:00")
	}
	if times_strings[1] != "10:00" {
		t.Fatalf(`formatTimesToString() converted 1000 to %s, expected %s.`, times_strings[1], "10:00")
	}
}

// TestBinarySearch | Test to see if binary search algorithm works correctly.
func TestBinarySearch(t *testing.T) {
	times := [...]int{
		800, 815, 830, 845, 900, 915, 930, 945, 1000, 1015, 1030,
		1100, 1115, 1130, 1145, 1200, 1215, 1230, 1245, 1300,
		1315, 1330, 1345, 1400, 1415, 1430, 1445, 1500, 1515, 1530,
		1545, 1600, 1615, 1630, 1645, 1700, 1715, 1730, 1745, 1800,
		1815, 1830, 1845, 1900, 1915, 1930, 1945, 2000, 2015, 2030,
		2045, 2100, 2115, 2130, 2145, 2200, 2215, 2230, 2245, 2300,
		2315, 2330, 2345,
	}

	expected_index := 4
	resulting_index := binary_search(times, 900)
	fmt.Println(resulting_index)

	if expected_index != resulting_index {
		t.Fatalf(`expected index to be %d but it was %d`, expected_index, resulting_index)
	}

}

func TestReturnTimeslotsInbetween(t *testing.T) {
	expected_result_range := [...]int{900, 915, 930, 945, 1000}

	start_time := 900
	end_time := 1000

	time_slots, err := return_time_slots_in_between(start_time, end_time)

	if err != nil {
		t.Fatalf(`TestReturn_time_slots_in_between failed completely with start_time: %d and end_time: %d`, start_time, end_time)
	}
	for index, time_slot := range *time_slots {
		if time_slot != expected_result_range[index] {
			t.Fatalf(`expected time slot to be %d but it was %d`, expected_result_range[index], time_slot)
		}
	}
}
