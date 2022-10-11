/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package raflaamoRestaurantsApi

import (
	"testing"
)

func TestRaflaamoRestaurantsApi_GetRestaurants(t *testing.T) {
	raflaamoRestaurantsApi := GetRaflaamoRestaurantsApi("rovaniemi")

	var currentTime int64 = 43200
	restaurantsFromApi, err := raflaamoRestaurantsApi.GetRestaurantsFromRaflaamoApi(currentTime)
	if err != nil {
		t.Errorf("[TestRaflaamoRestaurantsApi_GetRestaurants] (getRestaurantsFromRaflaamoApi) did not expect error but we got one: %s", err)
	}

	if len(restaurantsFromApi) == 0 {
		t.Errorf("[TestRaflaamoRestaurantsApi_GetRestaurants] (restaurantsFromApi) expected length of result to be < 0 but it was 0")
	}
}
