package main

import (
	"strconv"
	"strings"
	"testing"
)

// TestGetRestaurantsFromCity | Test to see if JSON parsing works correctly.
func TestGetRestaurantsFromCity(t *testing.T) {
	restaurant_data := getRestaurants()
	city := "helsinki"
	restaurants_from_helsinki := getRestaurantsFromCity(&*restaurant_data, &city)

	for _, restaurant := range *restaurants_from_helsinki {
		if !strings.Contains(strings.ToLower(*restaurant.Address.Municipality.Fi_FI), strings.ToLower(city)) {
			t.Errorf("restaurant.Address.Municipality.Fi_FI = %s, expected %s", *restaurant.Address.Municipality.Fi_FI, "helsinki")
		}
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

	times_strings := formatTimesToString(times_int)

	if times_strings[0] != "09:00" {
		t.Fatalf(`formatTimesToString() converted 900 to %s, expected %s.`, times_strings[0], "09:00")
	}
	if times_strings[1] != "10:00" {
		t.Fatalf(`formatTimesToString() converted 1000 to %s, expected %s.`, times_strings[1], "10:00")
	}
}
