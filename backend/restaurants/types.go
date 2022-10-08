/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package restaurants

import (
	"backend/raflaamoGraphApi"
	"backend/raflaamoRestaurantsApi"
	"backend/raflaamoTime"
	"regexp"
)

var regexToMatchRestaurantId = regexp.MustCompile(`[^fi/]\d+`)
var RegexToMatchTime = regexp.MustCompile(`\d{2}:\d{2}`)
var RegexToMatchDate = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)

type Restaurants struct {
	City                   string
	AmountOfEaters         int
	AllNeededRaflaamoTimes *raflaamoTime.RaflaamoTimes
	GraphApi               *raflaamoGraphApi.RaflaamoGraphApi
	RestaurantsApi         *raflaamoRestaurantsApi.RaflaamoRestaurantsApi
}
