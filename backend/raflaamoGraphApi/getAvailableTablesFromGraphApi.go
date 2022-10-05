package raflaamoGraphApi

import (
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

func (graphApi *RaflaamoGraphApi) getGraphApiResponseFromTimeSlot(requestUrlContainingTimeSlot string) (*parsedGraphData, error) {
	httpRequest := graphApi.getRaflaamoGraphApiRequest(requestUrlContainingTimeSlot)
	response, err := graphApi.sendRequestToGraphApi(httpRequest)
	if err != nil {
		return nil, err
	}
	deserializedGraphApiResponse, err := graphApi.deserializeGraphApiResponse(response)
	if err != nil {
		return nil, err
	}
	return deserializedGraphApiResponse, nil
}