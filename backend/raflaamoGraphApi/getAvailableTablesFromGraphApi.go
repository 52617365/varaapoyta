package raflaamoGraphApi

import (
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
	deserializedGraphApiResponse, err := deserializeGraphApiResponse(graphApi.response)
	if err != nil {
		return err
	}

	graphApi.deserializedResponse = deserializedGraphApiResponse
	return nil
}

// getAvailableTablesFromGraphApi should be called with a requestUrl payload that has already been initialized.
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
