package raflaamoGraphApi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// RaflaamoGraphApi control flow is: getRaflaamoGraphApiRequest -> interactWithGraphApi -> deserializeGraphApiResponse
type RaflaamoGraphApi struct {
	httpClient           *http.Client     // This will be initialized once because it does not change.
	request              *http.Request    // This will be initialized per request.
	response             *http.Response   // This will be initialized per request.
	deserializedResponse *parsedGraphData // This will be initialized per request.
}

func getRaflaamoGraphApi() *RaflaamoGraphApi {
	httpClient := &http.Client{}

	return &RaflaamoGraphApi{httpClient: httpClient}
}

func (graphApi *RaflaamoGraphApi) getRaflaamoGraphApiRequest(requestUrl string) {
	r, err := http.NewRequest("GET", requestUrl, nil)

	if err != nil {
		log.Fatal("[getRaflaamoGraphApiRequest] - Error initializing get request")
	}

	r.Header.Add("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	graphApi.request = r
}

func (graphApi *RaflaamoGraphApi) sendRequestToGraphApi() error {
	response, err := graphApi.httpClient.Do(graphApi.request)

	if err != nil {
		return fmt.Errorf("[interactWithApi] - %w", err)
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("[interactWithApi] - %w", errors.New("raflaamo api returned non 200 status code"))
	}

	graphApi.response = response
	return nil
}

func (graphApi *RaflaamoGraphApi) deserializeGraphApiResponse() error {
	response := graphApi.response

	var deserializedGraphData []parsedGraphData
	err := json.NewDecoder((response).Body).Decode(&deserializedGraphData)
	if err != nil {
		return fmt.Errorf("[deserializeGraphResponse] - %w", err)
	}
	if deserializedGraphData == nil {
		return errors.New("[deserializeGraphResponse] - there was an error deserializing the data")
	}

	// The relevant data is in the first index only.
	graphApi.deserializedResponse = &deserializedGraphData[0]
	return nil
}

// getAvailableTablesFromGraphApi should be called with a requestUrl payload that has already been initialized.
// TODO: get the time slots here.
//
// TODO: Also consider the restaurant kitchens closing time to avoid getting times where
//	 		 They don't let you reserve anymore (45 minutes before kitchen closes).
//
// TODO: don't handle time slots that are in the past from current time.
func (graphApi *RaflaamoGraphApi) getAvailableTablesFromGraphApi(requestUrl string) error {
	graphApi.getRaflaamoGraphApiRequest(requestUrl)

	var err error

	err = graphApi.sendRequestToGraphApi()
	if err != nil {
		return err
	}

	err = graphApi.deserializeGraphApiResponse()
	if err != nil {
		return err
	}
	return nil
}
