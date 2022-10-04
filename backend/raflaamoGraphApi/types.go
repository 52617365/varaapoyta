package raflaamoGraphApi

import (
	"backend/raflaamoRestaurantsApi"
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
	availableTimeIntervals chan *parsedGraphData
	restaurant             *ResponseFields
}
