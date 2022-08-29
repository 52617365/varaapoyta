package main

import (
	"fmt"
	"log"
	"testing"
)

func BenchmarkMain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		city := "helsinki"
		restaurants, _ := get_all_restaurants_from_raflaamo_api()

		if len(restaurants) == 0 {
			log.Fatalln("no restaurants found")
		}
		// here restaurants is not empty (we check it before)

		results := get_available_tables(city, restaurants, 1)
		for _, available_table := range results {
			start_string := fmt.Sprintf("name of restaurant: %s | available_tables: ", available_table.Name.Fi_FI)
			fmt.Println(start_string)

			for _, time := range available_table.available_time_slots {
				fmt.Println(time)
			}
		}
	}
}
