package main

import (
	"testing"
)

// TestGetRestaurants We expect response to be len(470).
func TestGetRestaurants(t *testing.T) {
	restaurants := getRestaurants()

	restaurants_length := len(*restaurants)
	if restaurants_length != 470 {
		t.Errorf("len(getRestaurants()) = %d, expected %d", restaurants_length, 470)
	}
}

func BenchmarkGetRestaurants(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getRestaurants()
	}
}
