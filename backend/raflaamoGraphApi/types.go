/*
 * Copyright (c) 2022. Rasmus Mäki
 */

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

type Response struct {
	availableTimeIntervals chan *graphApiResponseStructure.ParsedGraphData
	restaurant             *raflaamoRestaurantsApi.ResponseFields
}
