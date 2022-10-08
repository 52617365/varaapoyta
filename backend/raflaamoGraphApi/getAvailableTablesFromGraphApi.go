package raflaamoGraphApi

import (
	"backend/graphApiResponseStructure"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func GetRaflaamoGraphApi() *RaflaamoGraphApi {
	httpClient := &http.Client{}

	return &RaflaamoGraphApi{httpClient: httpClient}
}

func (graphApi *RaflaamoGraphApi) getRaflaamoGraphApiRequest(requestUrl string) *http.Request {
	r, err := http.NewRequest("GET", requestUrl, nil)

	if err != nil {
		// TODO: figure out if this can fail depending on requestUrl.
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
		return nil, fmt.Errorf("[sendRequestToGraphApi] - %w", errors.New("raflaamo graph api returned non 200 status code"))
	}

	return response, nil
}

func (graphApi *RaflaamoGraphApi) deserializeGraphApiResponse(graphApiResponse *http.Response) (*graphApiResponseStructure.ParsedGraphData, error) {
	var deserializedGraphData []graphApiResponseStructure.ParsedGraphData
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

func (graphApi *RaflaamoGraphApi) GetGraphApiResponseFromTimeSlot(requestUrlContainingTimeSlot string) (*graphApiResponseStructure.ParsedGraphData, error) {
	httpRequest := graphApi.getRaflaamoGraphApiRequest(requestUrlContainingTimeSlot)
	response, err := graphApi.sendRequestToGraphApi(httpRequest)
	if err != nil {
		return nil, err
	}
	deserializedGraphApiResponse, err := graphApi.deserializeGraphApiResponse(response)
	if err != nil {
		return nil, err
	}
	timeIntervals := *deserializedGraphApiResponse.Intervals
	if timeIntervals[0].Color == "transparent" {
		return nil, fmt.Errorf("[GetGraphResponseFromTimeSlot] - %w", errors.New("no available timeslots found"))
	}
	return deserializedGraphApiResponse, nil
}
