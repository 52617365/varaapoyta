package raflaamoGraphApi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type parsedGraphData struct {
	Name      string                `json:"name"`
	Intervals *[]parsedIntervalData `json:"intervals"` // were only interested in the first index.
	Id        int                   `json:"id"`
}

type parsedIntervalData struct {
	From  int64  `json:"from"`  // From is a unix timestamp in ms.
	To    int64  `json:"to"`    // To is a unix timestamp in ms.
	Color string `json:"color"` // Optional field, we can match this to see if the restaurant has available tables. (if not nil it does.)
}

func deserializeGraphApiResponse(res *http.Response) (*parsedGraphData, error) {
	var responseDecoded []parsedGraphData
	err := json.NewDecoder((res).Body).Decode(&responseDecoded)
	if err != nil {
		return nil, fmt.Errorf("[deserializeGraphResponse] - %w", err)
	}
	// Returning only the first index because the api for some reason contains weird data on top of the one we care about.
	// The relevant data is in the first index.
	if responseDecoded == nil {
		return nil, errors.New("[deserializeGraphResponse] - there was an error deserializing the data")
	}
	return &responseDecoded[0], nil
}
