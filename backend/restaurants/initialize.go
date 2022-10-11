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

func getInitializeProgram(city string, amountOfEaters string) *InitializeProgram {
	allNeededRaflaamoTimes := raflaamoTimes.GetAllNeededRaflaamoTimes()
	graphApi := raflaamoGraphApi.GetRaflaamoGraphApi()
	initializedRaflaamoRestaurantsApi := raflaamoRestaurantsApi.GetRaflaamoRestaurantsApi(city)
	return &InitializeProgram{
		City:                   city,
		AmountOfEaters:         amountOfEaters,
		AllNeededRaflaamoTimes: allNeededRaflaamoTimes,
		GraphApi:               graphApi,
		RestaurantsApi:         initializedRaflaamoRestaurantsApi,
	}
}
