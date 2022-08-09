package main

import (
	"testing"
)

// TestGetRestaurants We expect response to be len(470).
func TestGetRestaurants(t *testing.T) {
	restaurants := getAllRestaurantsFromRaflaamoApi()

	restaurants_length := len(*restaurants)
	if restaurants_length < 10 {
		// Can't check against a static number cuz the amount changes.
		t.Errorf("len(getRestaurants()) = %d, expected %s", restaurants_length, ">10")
	}
}

func BenchmarkGetRestaurants(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getAllRestaurantsFromRaflaamoApi()
	}
}
