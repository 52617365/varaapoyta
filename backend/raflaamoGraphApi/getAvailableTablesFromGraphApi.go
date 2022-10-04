package raflaamoGraphApi

import (
	"backend/raflaamoRestaurantsApi"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type ResponseFields = raflaamoRestaurantsApi.ResponseFields

// RaflaamoGraphApi control flow is: getRaflaamoGraphApiRequest -> interactWithGraphApi -> deserializeGraphApiResponse
type RaflaamoGraphApi struct {
	httpClient        *http.Client
	GraphApiResponses chan *Response
}

// Response TODO: Graph api responses should be stored with this struct.
type Response struct {
	availableTimeIntervals *parsedGraphData
	restaurant             *ResponseFields
}

func getRaflaamoGraphApi() *RaflaamoGraphApi {
	httpClient := &http.Client{}

	return &RaflaamoGraphApi{httpClient: httpClient}
}

func (graphApi *RaflaamoGraphApi) getRaflaamoGraphApiRequest(requestUrl string) *http.Request {
	r, err := http.NewRequest("GET", requestUrl, nil)

	if err != nil {
		log.Fatal("[getRaflaamoGraphApiRequest] - Error initializing get request")
	}

	r.Header.Add("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	return r
}

func (graphApi *RaflaamoGraphApi) sendRequestToGraphApi(graphApiRequest *http.Request) (*http.Response, error) {
	response, err := graphApi.httpClient.Do(graphApiRequest)

	if err != nil {
		return nil, fmt.Errorf("[sendRequestToGraphApi] - %w", err)
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("[sendRequestToGraphApi] - %w", errors.New("raflaamo api returned non 200 status code"))
	}

	return response, nil
}

func (graphApi *RaflaamoGraphApi) deserializeGraphApiResponse(graphApiResponse *http.Response) (*parsedGraphData, error) {
	var deserializedGraphData []parsedGraphData
	err := json.NewDecoder((graphApiResponse).Body).Decode(&deserializedGraphData)
	if err != nil {
		return nil, fmt.Errorf("[deserializeGraphApiResponse] - %w", err)
	}
	if deserializedGraphData == nil {
		return nil, errors.New("[deserializeGraphApiResponse] - there was an error deserializing the data")
	}
	// The relevant data from the graph API is in the first index only.
	return &deserializedGraphData[0], nil
}

// getAvailableTablesFromGraphApi should be called with a requestUrl payload that has already been initialized.
// TODO: get the time slots here.
// TODO: Also consider the restaurant kitchens closing time to avoid getting times when they don't let you reserve anymore (45 minutes before kitchen closes).
// TODO: don't handle time slots that are in the past from current time.
// TODO: Use channels and goroutines.

// Idea is to have something iterating the restaurants that calls this function.
// The caller should be the one initializing the Response struct and the available tables just get assigned to it here.
// The restaurant gets assigned in the callee function and this function does not have to worry about that.
func (graphApi *RaflaamoGraphApi) getAvailableTablesFromGraphApi(requestUrl string, graphApiResponse *Response) error {
	httpRequest := graphApi.getRaflaamoGraphApiRequest(requestUrl)

	response, err := graphApi.sendRequestToGraphApi(httpRequest)
	if err != nil {
		return err
	}

	deserializedGraphApiResponse, err := graphApi.deserializeGraphApiResponse(response)
	if err != nil {
		return err
	}

	graphApiResponse.availableTimeIntervals = deserializedGraphApiResponse
	return nil
}
