package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type parsed_graph_data struct {
	Name      string                 `json:"name"`
	Intervals []parsed_interval_data `json:"intervals"` // were only interested in the first index.
	Id        int                    `json:"id"`
}

type parsed_interval_data struct {
	From  int64  `json:"from"`  // From is a unix timestamp in ms.
	To    int64  `json:"to"`    // To is a unix timestamp in ms.
	Color string `json:"color"` // Optional field, we can match this to see if the restaurant has available tables. (if not nil it does.)
}

func deserialize_graph_response(res **http.Response) (*parsed_graph_data, error) {
	var response_decoded []parsed_graph_data
	err := json.NewDecoder((*res).Body).Decode(&response_decoded)
	if err != nil {
		return nil, err
	}
	// Returning only the first index because the api for some reason contains weird data on top of the one we care about.
	// The relevant data is in the first index.
	if response_decoded == nil {
		return nil, errors.New("there was an error deserializing the data")
	}
	return &response_decoded[0], nil
}
