/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package restaurants

import (
	"testing"
)


func BenchmarkGetRestaurantsAndAvailableTables(b *testing.B) {
	for i := 0; i < b.N; i++ {
		restaurantsInstance, _ := GetRestaurants("helsinki", "1")
		raflaamoRestaurants := restaurantsInstance.GetRestaurantsAndAvailableTables()
		for _, restaurant := range raflaamoRestaurants {
			if <-restaurant.GraphApiResults.Err != nil {
				continue
			}
			for range restaurant.GraphApiResults.AvailableTimeSlotsBuffer {
			}
		}
	}
}
