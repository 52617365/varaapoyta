package main

import (
	"fmt"
	"log"
	"testing"
)

func BenchmarkMain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		city := "helsinki"

		init_request, _ := init_restaurants()
		response, err := init_request.get()

		if len(response) == 0 && err == nil {
			log.Fatalln("expected error")
		}
		// here restaurants is not empty (we check it before)
		results, _ := get_available_tables(city, 1)
		for _, available_table := range results {
			start_string := fmt.Sprintf("name of restaurant: %s | available_tables: ", available_table.Name.Fi_FI)
			fmt.Println(start_string)

			for _, time := range available_table.Available_time_slots {
				fmt.Println(time)
			}
		}
	}
}
