package main

import (
	"errors"
	"fmt"
	"net/http"
)

// Gets timeslots from raflaamo API that is responsible for returning graph data.
// Instead of drawing a graph with it, we convert it into time to determine which table is open or not.
func get_time_slots_from_graph_api(id_from_reservation_page_url string, current_date string, time_slot string, amount_of_eaters int) (*parsed_graph_data, error) {
	// https://s-varaukset.fi/api/recommendations/slot/{id}/{date}/{time}/{amount_of_eaters}
	request_url := fmt.Sprintf("https://s-varaukset.fi/api/recommendations/slot/%s/%s/%s/%d", id_from_reservation_page_url, current_date, time_slot, amount_of_eaters)

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
	deserialized_graph_data_with_open_tables := filter_out_time_slots_with_no_open_tables(deserialized_graph_data)

	if deserialized_graph_data_with_open_tables == nil {
		return nil, errors.New("there were no open tables")
	}
	return deserialized_graph_data_with_open_tables, nil
}

func filter_out_time_slots_with_no_open_tables(datas []parsed_graph_data) *parsed_graph_data {
	times_with_open_tables := make([]parsed_graph_data, 0, len(datas))
	for _, data := range datas {
		if !time_slot_does_not_contain_open_tables(data) {
			times_with_open_tables = append(times_with_open_tables, data)
		}
	}
	// Returning only the first index because the api for some reason contains weird data on top of the one we care about.
	return &times_with_open_tables[0]
}
