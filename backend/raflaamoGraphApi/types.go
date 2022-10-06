package raflaamoGraphApi

import (
	"backend/raflaamoRestaurantsApi"
	"backend/raflaamoTime"
	"net/http"
)

type ResponseFields = raflaamoRestaurantsApi.ResponseFields
type RaflaamoTimes = raflaamoTime.RaflaamoTimes

// RaflaamoGraphApi control flow is: getRaflaamoGraphApiRequest -> InteractWithGraphApi -> DeserializeGraphApiResponse
type RaflaamoGraphApi struct {
	httpClient        *http.Client
	GraphApiResponses chan *Response
}

// Response TODO: Graph api responses should be stored with this struct.
type Response struct {
	availableTimeIntervals chan *ParsedGraphData
	restaurant             *ResponseFields
}
