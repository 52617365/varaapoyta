package restaurants

import "testing"

// ~20 seconds pre goroutines.
func BenchmarkGetRestaurantsAndAvailableTables(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetRestaurantsAndAvailableTables("helsinki", 1)
	}
}
