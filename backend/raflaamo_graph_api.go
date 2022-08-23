package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/exp/slices"
)

// Gets timeslots from raflaamo API that is responsible for returning graph data.
// In the end, instead of drawing a graph with it, we convert it into time to determine which table is open or not.
// This one sends requests so we use goroutines in it.
// TODO: figure out how we should use goroutines and channels here whilst keeping it simple.
func get_available_time_intervals_from_graph_api(id_from_reservation_page_url string, all_possible_time_slots []covered_times, amount_of_eaters int, all_reservation_times []int64, current_time date_and_time) ([]string, error) {
	// Reserve space for all the time_slots which the function will return in the end.
	time_slots_from_graph_api := make([]string, 0, len(all_possible_time_slots))

	results := make(chan string)
	// Looping through all the poss
	for _, possible_time_slot := range all_possible_time_slots {
		possible_time_slot := possible_time_slot
		go func() {
			err := func() error {
				time_slot_string := get_string_time_from_unix(possible_time_slot.time)
				// replacing the 17(:)00 to match the format in url.
				time_slot_string = strings.Replace(time_slot_string, ":", "", -1)
				// Example of a request_url: https://s-varaukset.fi/api/recommendations/slot/{id}/{date}/{time}/{amount_of_eaters}
				request_url := fmt.Sprintf("https://s-varaukset.fi/api/recommendations/slot/%s/%s/%s/%d", id_from_reservation_page_url, current_time.date, time_slot_string, amount_of_eaters)

				r, err := http.NewRequest("GET", request_url, nil)

				if err != nil {

					return errors.New("error connecting to raflaamo api")
				}
				r.Header.Add("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

				client := &http.Client{}
				res, err := client.Do(r)

				// Will throw if we call deserialize_graph_response with a status code other than 200, so we handle it here.
				if err != nil || res.StatusCode != 200 {
					return errors.New("error connecting to raflaamo api")
				}

				// gets stuck here for some reason?
				deserialized_graph_data, err := deserialize_graph_response(&res)

				// most likely won't jump into this branch but check regardless.
				if err != nil {
					return errors.New("there was an error deserializing the data returned from endpoint")
				}
				if time_slot_does_not_contain_open_tables(*deserialized_graph_data) {
					return errors.New("time slot did not contain open tables")
				}

				// Adding 10800000(ms) to the time to match utc +2 or +3 (finnish time) (10800000 ms corresponds to 3h)
				// because graph unix time fields "to" and "from" come in utc+0
				graph_end_unix := deserialized_graph_data.Intervals[0].To
				graph_end_unix += 10800000

				time_slots, err := time_slots_in_between(current_time.time, graph_end_unix, all_reservation_times)
				if err != nil {
					return errors.New("error in time_slots_in_between")
				}
				// Avoiding storing duplicate time slots because without this, it will.
				for _, time_slot := range time_slots {
					results <- time_slot
				}
				return nil
			}()
			if err != nil {
				return
			}
		}()
	}
	for _ = range results {
		if !slices.Contains(time_slots_from_graph_api, <-results) {
			time_slots_from_graph_api = append(time_slots_from_graph_api, <-results)
		}
	}
	return time_slots_from_graph_api, nil
}
