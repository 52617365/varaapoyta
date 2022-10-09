/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package raflaamoRestaurantsApi

import "net/http"

type RaflaamoRestaurantsApi struct {
	httpClient               *http.Client
	request                  *http.Request
	response                 *http.Response
	cityToGetRestaurantsFrom string
	currentTime              int64
}
