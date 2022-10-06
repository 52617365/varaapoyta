package raflaamoGraphApi

import (
	"backend/graphApiResponseStructure"
	"backend/raflaamoRestaurantsApi"
	"net/http"
)

// RaflaamoGraphApi control flow is: getRaflaamoGraphApiRequest -> InteractWithGraphApi -> DeserializeGraphApiResponse
type RaflaamoGraphApi struct {
	httpClient        *http.Client
	GraphApiResponses chan *Response
}

// Response TODO: Graph api responses should be stored with this struct.
type Response struct {
	availableTimeIntervals chan *graphApiResponseStructure.ParsedGraphData
	restaurant             *raflaamoRestaurantsApi.ResponseFields
}
