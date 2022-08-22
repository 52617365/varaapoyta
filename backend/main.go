package main

import (
	"fmt"
	"log"
)

// TODO: Make endpoints.
func main() {
	restaurants, err := filter_restaurants_from_city("rovaniemi")
	if err != nil {
		// if error we return this from the endpoint.
		log.Fatalln(err)
	}
	if len(restaurants) == 0 {
		log.Fatalln("no restaurants found")
	}
	// here restaurants is not empty (we check it before)

	current_time := get_current_date_and_time()
	available_tables := get_available_tables(restaurants, current_time, 1)
	for _, available_table := range available_tables {
		start_string := fmt.Sprintf("name of restaurant: %s | available_tables: ", available_table.restaurant.Name.Fi_FI)
		fmt.Println(start_string)

		for _, time := range available_table.available_time_slots {
			fmt.Println(time)
		}
	}
}
