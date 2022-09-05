package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gammazero/workerpool"
	"golang.org/x/exp/slices"
)

func construct_payload(id_from_reservation_page_url string, current_date string, time covered_times, amount_of_eaters int) string {
	time_slot_string := get_string_time_from_unix(time.time)
	// replacing the 17(:)00 to match the format in url.
	time_slot_string = strings.Replace(time_slot_string, ":", "", -1)
	// Example of a request_url: https://s-varaukset.fi/api/recommendations/slot/{id}/{date}/{time}/{amount_of_eaters}
	request_url := fmt.Sprintf("https://s-varaukset.fi/api/recommendations/slot/%s/%s/%s/%d", id_from_reservation_page_url, current_date, time_slot_string, amount_of_eaters)
	return request_url
}

// Gets timeslots from raflaamo API that is responsible for returning graph data.
// In the end, instead of drawing a graph with it, we convert it into time to determine which table is open or not.
// This one sends requests, so we use goroutines in it.
// TODO: this function is our bottleneck. We have to do something about this.
func interact_with_api(all_possible_time_slots []covered_times, id_from_reservation_page_url string, current_date string, amount_of_eaters int) (chan *parsed_graph_data, error) {
	wp := workerpool.New(5)
	response_chan := make(chan *parsed_graph_data, len(all_possible_time_slots))

	for _, time := range all_possible_time_slots {
		time := time
		wp.Submit(func() {
			request_url := construct_payload(id_from_reservation_page_url, current_date, time, amount_of_eaters)

			r, err := http.NewRequest("GET", request_url, nil)

			if err != nil {
				return
			}

			r.Header.Add("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
			client := &http.Client{}
			res, err := client.Do(r)

			// Will throw if we call deserialize_graph_response with a status code other than 200, so we handle it here.
			if err != nil || res.StatusCode != 200 {
				return
			}

			deserialized_graph_data, err := deserialize_graph_response(&res)
			// most likely won't jump into this branch but check regardless.
			if err != nil {
				return
			}
			if time_slot_does_not_contain_open_tables(*deserialized_graph_data) {
				return
			}
			// Adding timezone difference into the unix time. (three hours).
			graph_end_unix := deserialized_graph_data.Intervals[0].To
			// Adding 10800000(ms) to the time to match utc +2 or +3 (finnish time) (10800000 ms corresponds to 3h)
			// because graph unix time fields "to" and "from" come in utc+0
			graph_end_unix += 10800000
			response_chan <- deserialized_graph_data
		})
	}
	wp.StopWait()
	close(response_chan)
	return response_chan, nil
}

func get_available_time_intervals_from_graph_api(restaurant_starting_time_unix int64, restaurant_closing_time_unix int64, id_from_reservation_page_url string, time_slots_to_check_from_graph_api []covered_times, amount_of_eaters int, all_reservation_times []int64, current_time date_and_time) ([]string, error) {
	// Reserve space for all the time_slots which the function will return in the end.
	api_responses, err := interact_with_api(time_slots_to_check_from_graph_api, id_from_reservation_page_url, current_time.date, amount_of_eaters)

	// most likely won't jump into this branch but check regardless.
	if err != nil {
		return nil, err
	}

	time_slots := make([]string, 0, len(time_slots_to_check_from_graph_api))
	for api_response := range api_responses {
		graph_end_unix := api_response.Intervals[0].To

		if current_time.time >= graph_end_unix {
			return nil, errors.New("trying to get a time_slot with invalid timestamps")
		}

		// Avoiding storing duplicate time slots because without this, it will.
		for _, reservation_time := range all_reservation_times {
			// We check if the timestamps are valid here.
			if valid_graph_time_slot(reservation_time, current_time.time, graph_end_unix) && time_slot_in_restaurant_opening_hours(reservation_time, restaurant_starting_time_unix, restaurant_closing_time_unix) {
				slot := get_string_time_from_unix(reservation_time)
				if !slices.Contains(time_slots, slot) {
					time_slots = append(time_slots, slot)
				}
			}
		}
	}
	return time_slots, nil
}

// Checks to see if the reservation time checked is larger than the current time and smaller or equal to the last possible time to reserve (graph_end_unix)
func valid_graph_time_slot(reservation_time int64, current_time int64, graph_end_unix int64) bool {
	if reservation_time > current_time && reservation_time <= graph_end_unix {
		return true
	}
	return false
}

func time_slot_in_restaurant_opening_hours(reservation_time int64, restaurant_starting_time_unix int64, restaurant_closing_time_unix int64) bool {
	if reservation_time > restaurant_starting_time_unix && reservation_time <= restaurant_closing_time_unix {
		return true
	}
	return false
}
