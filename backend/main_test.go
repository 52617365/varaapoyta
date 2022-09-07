package main

import (
	"fmt"
	"log"
	"testing"
)

func BenchmarkMain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		city := "helsinki"
		raflaamo_api_response_chan := make(chan []response_fields)
		raflaamo_api_response_error_chan := make(chan error)
		go get_all_restaurants_from_raflaamo_api(raflaamo_api_response_chan, raflaamo_api_response_error_chan)

		raflaamo_api_response := <-raflaamo_api_response_chan
		raflaamo_api_response_error := <-raflaamo_api_response_error_chan
		if len(raflaamo_api_response) == 0 && raflaamo_api_response_error == nil {
			log.Fatalln("expected error")
		}
		// here restaurants is not empty (we check it before)
		results := get_available_tables(city, 1)
		for _, available_table := range results {
			start_string := fmt.Sprintf("name of restaurant: %s | available_tables: ", available_table.Name.Fi_FI)
			fmt.Println(start_string)

			for _, time := range available_table.Available_time_slots {
				fmt.Println(time)
			}
		}
	}
}
