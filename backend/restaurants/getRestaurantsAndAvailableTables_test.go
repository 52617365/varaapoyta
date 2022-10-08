package restaurants

import (
	"log"
	"testing"
)

// ~20 seconds pre goroutines.
func BenchmarkGetRestaurantsAndAvailableTables(b *testing.B) {
	for i := 0; i < b.N; i++ {
		restaurantsInstance, _ := GetRestaurants("helsinki", 1)
		raflaamoRestaurants, err := restaurantsInstance.GetRestaurantsAndAvailableTables()
		if err != nil {
			log.Fatalln("err")
		}
		for _, restaurant := range raflaamoRestaurants {
			if <-restaurant.GraphApiResults.Err != nil {
				continue
			}
			for range restaurant.GraphApiResults.AvailableTimeSlotsBuffer {
			}
		}
	}
}
