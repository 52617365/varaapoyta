package main

import (
	"fmt"
	"log"
	"testing"
)

func BenchmarkMain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		city := "helsinki"

		initRequest, _ := init_restaurants()
		response, err := initRequest.get()

		if len(response) == 0 && err == nil {
			log.Fatalln("expected error")
		}
		// here restaurants is not empty (we check it before)
		results, _ := get_available_tables(city, 1)
		for _, availableTable := range results {
			startString := fmt.Sprintf("name of restaurant: %s | available_tables: ", availableTable.Name.Fi_FI)
			fmt.Println(startString)

			for _, time := range availableTable.Available_time_slots {
				fmt.Println(time)
			}
		}
	}
}
