/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package restaurants

import (
	"backend/raflaamoGraphApi"
	"backend/raflaamoRestaurantsApi"
	"backend/raflaamoTimes"
)

type InitializeProgram struct {
	City                   string
	AmountOfEaters         string
	AllNeededRaflaamoTimes *raflaamoTimes.RaflaamoTimes
	GraphApi               *raflaamoGraphApi.RaflaamoGraphApi
	RestaurantsApi         *raflaamoRestaurantsApi.RaflaamoRestaurantsApi
}
