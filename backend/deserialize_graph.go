package main

import (
	"encoding/json"
	"net/http"
)

type parsed_graph_data struct {
	Name      string                 `json:"name"`
	Intervals []parsed_interval_data `json:"intervals"` // were only interested in the first index.
	Seats     int                    `json:"seats"`
	Id        int                    `json:"id"`
	Env       interface{}            `json:"env:"`
	Host      string                 `json:"host"`
}

type parsed_interval_data struct {
	From  int     `json:"from"`  // From is a unix timestamp in ms.
	To    int     `json:"to"`    // To is a unix timestamp in ms.
	Color *string `json:"color"` // Optional field, we can match this to see if the restaurant has available tables. (if not nil it does.)
}

func deserialize_graph_response(res **http.Response) *parsed_graph_data {
	var response_decoded []parsed_graph_data
	err := json.NewDecoder((*res).Body).Decode(&response_decoded)
	if err != nil {
		return nil
	}
	return &response_decoded[0]
}
