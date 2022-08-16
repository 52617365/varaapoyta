package main

import (
	"fmt"
	"log"
)

// TODO: Make endpoints.
// TODO: Restaurant often don't take reservations in the 1h time slot before they close. Check the closing time and don't include reservation times that are in the 1h window before closing.

func main() {
	city := "helsinki"
	amount_of_eaters := 1
	restaurants, err := filter_restaurants_from_city(city)
	if err != nil {
		log.Fatal("Could not find any restaurants.")
	}
	results := getAvailableTables(restaurants, amount_of_eaters)
	for _, result := range results {
		for _, time_slot := range result.available_time_slots {
			fmt.Println(time_slot)
		}
	}
}

// type restaurant_with_available_times_struct struct {
// 	restaurant           response_fields
// 	available_time_slots []string
// }
