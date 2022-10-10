/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package restaurants

import (
	"backend/raflaamoGraphApi"
	"backend/raflaamoRestaurantsApi"
	"backend/raflaamoTime"
)

type Restaurants struct {
	City                   string
	AmountOfEaters         string
	AllNeededRaflaamoTimes *raflaamoTime.RaflaamoTimes
	GraphApi               *raflaamoGraphApi.RaflaamoGraphApi
	RestaurantsApi         *raflaamoRestaurantsApi.RaflaamoRestaurantsApi
}
