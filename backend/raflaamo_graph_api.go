package main

import (
	"errors"
	"fmt"
	"golang.org/x/exp/slices"
	"net/http"
	"strings"
)

// Gets timeslots from raflaamo API that is responsible for returning graph data.
// In the end, instead of drawing a graph with it, we convert it into time to determine which table is open or not.
func get_available_time_intervals_from_graph_api(id_from_reservation_page_url string, all_possible_time_slots []covered_times, amount_of_eaters int, all_reservation_times []int64, current_time date_and_time) ([]string, error) {
	// Reserve space for all the time_slots which the function will return in the end.
	time_slots_from_graph_api := make([]string, 0, len(all_possible_time_slots))

	// Looping through all the poss
	for _, possible_time_slot := range all_possible_time_slots {
		//time_slot.time
		time_slot_string := get_string_time_from_unix(possible_time_slot.time)
		// replacing the 17(:)00 to match the format in url.
		time_slot_string = strings.Replace(time_slot_string, ":", "", -1)
		// https://s-varaukset.fi/api/recommendations/slot/{id}/{date}/{time}/{amount_of_eaters}
		request_url := fmt.Sprintf("https://s-varaukset.fi/api/recommendations/slot/%s/%s/%s/%d", id_from_reservation_page_url, current_time.date, time_slot_string, amount_of_eaters)

		r, err := http.NewRequest("GET", request_url, nil)

		if err != nil {
			return nil, errors.New("error connecting to raflaamo api")
		}
		r.Header.Add("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

		client := &http.Client{}
		res, err := client.Do(r)

		// Will throw if we call deserialize_graph_response with a status code other than 200, so we handle it here.
		if err != nil || res.StatusCode != 200 {
			return nil, errors.New("error connecting to raflaamo api")
		}

		deserialized_graph_data := deserialize_graph_response(&res)

		// most likely won't jump into this branch but check regardless.
		if deserialized_graph_data == nil {
			return nil, errors.New("there was an error deserializing the data returned from endpoint")
		}
		if time_slot_does_not_contain_open_tables(*deserialized_graph_data) {
			return nil, errors.New("time slot did not contain open tables")
		}

		// Adding 10800000(ms) to the time to match utc +2 or +3 (finnish time) (10800000 ms corresponds to 3h)
		// because graph unix times come in utc+0
		graph_end_unix := deserialized_graph_data.Intervals[0].To
		graph_end_unix += 10800000

		time_slots, err := time_slots_in_between(current_time.time, graph_end_unix, all_reservation_times)
		if err != nil {
			continue
		}
		// Avoiding storing duplicate time slots.
		for _, time_slot := range time_slots {
			if !slices.Contains(time_slots_from_graph_api, time_slot) {
				time_slots_from_graph_api = append(time_slots_from_graph_api, time_slot)
			}
		}
	}
	return time_slots_from_graph_api, nil
}

// for _, time_slot := range all_possible_time_slots {
// 	// TODO: store result in channel. It's 22/08/2022 and today I'm not big brained enough to do it.
// 	time_slots_from_graph_api, err := get_time_slots_from_graph_api(id_from_reservation_page_url, all_possible_time_slots, amount_of_eaters)
// 	if err != nil {
// 		continue // it's err if there was an error connecting to raflaamo API or if there were no results.
// 	}
// 	// time_slots_from_graph_api.Intervals[0].From += 10800000

// 	time_slots_from_graph_api.Intervals[0].To += 10800000
// 	graph_end_unix := time_slots_from_graph_api.Intervals[0].To
// 	time_slots, err := time_slots_in_between(current_time.time, graph_end_unix, all_reservation_times)
// 	if err != nil {
// 		continue
// 	}
